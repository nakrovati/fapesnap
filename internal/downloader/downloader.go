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

	"github.com/egoist/mygo"
	"github.com/nakrovati/fapesnap/internal/config"
	"github.com/nakrovati/fapesnap/internal/pkg/utils"
	"github.com/nakrovati/fapesnap/internal/providers"
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

//nolint:gochecknoglobals // package-level event registry is required for external subscriptions.
var (
	DownloadBatchStarted   = mygo.NewEvent[struct{}]("download:batch-started")
	DownloadBatchCompleted = mygo.NewEvent[DownloadBatchEvent]("download:batch-completed")
	DownloadMediaStarted   = mygo.NewEvent[DownloadMediaEvent]("download:media-started")
	DownloadMediaFailed    = mygo.NewEvent[DownloadMediaEvent]("download:media-failed")
	DownloadMediaCompleted = mygo.NewEvent[DownloadMediaEvent]("download:media-completed")
)

type Downloader struct {
	httpClient *http.Client
	logger     *slog.Logger
}

type MediaURLResolver func(pageURL string) (string, error)

func NewDownloader(logger *slog.Logger) *Downloader {
	if logger == nil {
		logger = slog.Default()
	}

	return &Downloader{
		logger:     logger,
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

	emitEvent(d.logger, DownloadBatchStarted, struct{}{})

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

				emitEvent(d.logger, DownloadMediaFailed, DownloadMediaEvent{
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

	emitEvent(d.logger, DownloadBatchCompleted, DownloadBatchEvent{
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
	emitEvent(d.logger, DownloadMediaStarted, DownloadMediaEvent{
		PageURL: pageURL,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		emitEvent(d.logger, DownloadMediaFailed, DownloadMediaEvent{
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
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		emitEvent(d.logger, DownloadMediaFailed, DownloadMediaEvent{
			PageURL: pageURL,
			Error:   err.Error(),
		})

		return fmt.Errorf("http request failed: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var statusErr error

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusForbidden:
		statusErr = fmt.Errorf("forbidden (headers/cookies/hotlink): %d, %s", resp.StatusCode, url)
	case http.StatusNotFound:
		statusErr = fmt.Errorf("media not found: %d, %s", resp.StatusCode, url)
	default:
		statusErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if statusErr != nil {
		emitEvent(d.logger, DownloadMediaFailed, DownloadMediaEvent{
			PageURL: pageURL,
			Error:   statusErr.Error(),
		})

		return statusErr
	}

	if err := d.SaveMedia(resp, url, dir); err != nil {
		emitEvent(d.logger, DownloadMediaFailed, DownloadMediaEvent{
			PageURL: pageURL,
			Error:   err.Error(),
		})

		return fmt.Errorf("save media: %w", err)
	}

	emitEvent(d.logger, DownloadMediaCompleted, DownloadMediaEvent{
		PageURL: pageURL,
	})

	return nil
}

func emitEvent[T any](logger *slog.Logger, event *mygo.Event[T], payload T) {
	if err := event.Broadcast(payload); err != nil {
		logger.Error("Failed to emit event", "event", event.Name(), "error", err)
	}
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
