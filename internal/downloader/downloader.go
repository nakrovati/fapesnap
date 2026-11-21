package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nakrovati/fapesnap/internal/config"
	"github.com/nakrovati/fapesnap/internal/pkg/utils"
	"github.com/nakrovati/fapesnap/internal/providers"
	"github.com/wailsapp/wails/v3/pkg/application"
)

var (
	ErrMediaNotFound = errors.New("media not found")
	ErrUsernameEmpty = errors.New("username cannot be empty")
)

type DownloadMediaEvent struct {
	PageURL string `json:"pageUrl"`
	Error   string `json:"error,omitempty"`
}

type DownloadBatchEvent struct {
	Total      int `json:"total"`
	Downloaded int `json:"downloaded,omitempty"`
	Failed     int `json:"failed,omitempty"`
}

type Downloader struct {
	httpClient *http.Client
	logger     *slog.Logger
	event      *application.EventManager
}

type MediaURLResolver func(pageURL string) (string, error)

func NewDownloader(logger *slog.Logger, event *application.EventManager) *Downloader {
	return &Downloader{
		logger:     logger,
		event:      event,
		httpClient: newHTTPClient(),
	}
}

func newHTTPClient() *http.Client {
	transport := http.DefaultTransport

	defaultTransport, ok := http.DefaultTransport.(*http.Transport)
	if ok {
		clonedTransport := defaultTransport.Clone()
		clonedTransport.ResponseHeaderTimeout = 60 * time.Second
		transport = clonedTransport
	}

	return &http.Client{
		Transport:     transport,
		CheckRedirect: nil,
		Jar:           nil,
		// No overall timeout: large files may take longer to transfer.
		Timeout: 0,
	}
}

func (d *Downloader) DownloadMediaItems(
	ctx context.Context,
	mediaItems []providers.Media,
	baseDownloadDir config.DownloadDir,
	providerName string,
	collectionSlug string,
	maxParallelDownloads int,
	resolveMediaURL MediaURLResolver,
) error {
	downloadDir, err := utils.GetCollectionDownloadDir(baseDownloadDir, providerName, collectionSlug)
	if err != nil {
		return fmt.Errorf("failed to get download directory: %w", err)
	}

	d.event.Emit("download:batch-started")

	jobs := make(chan providers.Media)

	var (
		wg                   sync.WaitGroup
		downloadedMediaCount int64
		failedMediaCount     int64
	)

	worker := func() {
		for media := range jobs {
			mediaURL, resolveErr := resolveMediaURLForItem(media, resolveMediaURL)
			if resolveErr != nil {
				atomic.AddInt64(&failedMediaCount, 1)

				d.logger.Error("Failed to resolve media URL", "pageUrl", media.PageURL, "error", resolveErr)

				d.event.Emit("download:media-failed", DownloadMediaEvent{
					PageURL: media.PageURL,
					Error:   resolveErr.Error(),
				})

				continue
			}

			err := d.DownloadMedia(ctx, media.PageURL, mediaURL, downloadDir)
			if err != nil {
				atomic.AddInt64(&failedMediaCount, 1)

				d.logger.Error("Failed to download media", "url", mediaURL, "error", err)

				continue
			}

			d.logger.Info("Media downloaded", "url", mediaURL)

			atomic.AddInt64(&downloadedMediaCount, 1)
		}
	}

	for range maxParallelDownloads {
		wg.Go(worker)
	}

	go func() {
		defer close(jobs)

		for _, media := range mediaItems {
			select {
			case <-ctx.Done():
				return
			case jobs <- media:
			}
		}
	}()

	wg.Wait()

	d.event.Emit("download:batch-completed", DownloadBatchEvent{
		Total:      len(mediaItems),
		Downloaded: int(downloadedMediaCount),
		Failed:     int(failedMediaCount),
	})

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}

func resolveMediaURLForItem(media providers.Media, resolver MediaURLResolver) (string, error) {
	if media.URL != "" {
		return media.URL, nil
	}

	if resolver == nil {
		return "", errors.New("media URL resolver is not configured")
	}

	return resolver(media.PageURL)
}

func (d *Downloader) DownloadMedia(ctx context.Context, pageURL, url string, dir string) error {
	d.event.Emit("download:media-started", DownloadMediaEvent{
		PageURL: pageURL,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		d.event.Emit("download:media-failed", DownloadMediaEvent{
			PageURL: pageURL,
			Error:   err.Error(),
		})

		return fmt.Errorf("build request: %w", err)
	}

	req.Header = http.Header{
		"User-Agent": {
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148 Safari/537.36",
		},
		"Accept":          {"image/avif,image/webp,image/apng,image/*,*/*;q=0.8"},
		"Accept-Language": {"en-US,en;q=0.9"},
		"Referer":         {deriveReferer(url)},
		"Connection":      {"keep-alive"},
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		d.event.Emit("download:media-failed", DownloadMediaEvent{
			PageURL: pageURL,
			Error:   err.Error(),
		})

		return fmt.Errorf("http request failed: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		// ok
	case http.StatusForbidden:
		return fmt.Errorf("forbidden (headers/cookies/hotlink): %d, %s", resp.StatusCode, url)
	case http.StatusNotFound:
		return fmt.Errorf("media not found: %d, %s", resp.StatusCode, url)
	default:
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := d.SaveMedia(resp, url, dir); err != nil {
		d.event.Emit("download:media-failed", DownloadMediaEvent{
			PageURL: pageURL,
			Error:   err.Error(),
		})

		return fmt.Errorf("save media: %w", err)
	}

	d.event.Emit("download:media-completed", DownloadMediaEvent{
		PageURL: pageURL,
	})

	return nil
}

func (d *Downloader) SaveMedia(resp *http.Response, src string, dir string) error {
	fileName := mediaFileName(src)
	filePath := filepath.Join(dir, fileName)

	cleanDir := filepath.Clean(dir)
	cleanPath := filepath.Clean(filePath)

	if !strings.HasPrefix(cleanPath, cleanDir) {
		return fmt.Errorf("file path escapes target directory: %s", filePath)
	}

	file, err := os.Create(cleanPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			fmt.Printf("Failed to close file: %v\n", err)
		}
	}()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func mediaFileName(src string) string {
	parsedURL, err := url.Parse(src)
	if err == nil {
		// Bunkr passes the original file name in the `n` query parameter.
		if name := strings.TrimSpace(parsedURL.Query().Get("n")); name != "" {
			name = filepath.Base(filepath.FromSlash(name))
			if name != "." && name != string(filepath.Separator) && name != "" {
				return sanitizeFileName(name)
			}
		}

		return filepath.Base(parsedURL.Path)
	}

	return filepath.Base(strings.Split(src, "?")[0])
}

func sanitizeFileName(name string) string {
	const invalidChars = `<>:"/\|?*`

	return strings.Map(func(r rune) rune {
		if r < 32 || strings.ContainsRune(invalidChars, r) {
			return '_'
		}

		return r
	}, name)
}

func deriveReferer(src string) string {
	u, err := url.Parse(src)
	if err != nil {
		return "https://example.com/"
	}

	return u.Scheme + "://" + u.Host + "/"
}
