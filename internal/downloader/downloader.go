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

type Downloader struct {
	httpClient *http.Client
	logger     *slog.Logger
	event      *application.EventManager
}

func NewDownloader(logger *slog.Logger, event *application.EventManager) *Downloader {
	return &Downloader{
		logger: logger,
		event:  event,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (d *Downloader) DownloadMediaItems(
	ctx context.Context,
	mediaItems []providers.Media,
	baseDownloadDir config.DownloadDir,
	providerName string,
	collectionSlug string,
	maxParallelDownloads int,
) error {
	downloadDir, err := utils.GetCollectionDownloadDir(baseDownloadDir, providerName, collectionSlug)
	if err != nil {
		return fmt.Errorf("failed to get download directory: %w", err)
	}

	d.event.Emit("download:started")

	jobs := make(chan providers.Media)

	var wg sync.WaitGroup
	var downloadedMediaCount int64

	worker := func() {
		for media := range jobs {
			err := d.DownloadMedia(ctx, media.URL, downloadDir)
			if err != nil {
				d.logger.Error("Failed to download media", "url", media.URL, "error", err)

				continue
			}

			d.logger.Info("Media downloaded", "url", media.URL)

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

	totalMediaCount := int64(len(mediaItems))

	if totalMediaCount > downloadedMediaCount {
		d.event.Emit("download:completed", fmt.Sprintf("%d files out of %d downloaded", downloadedMediaCount, len(mediaItems)))
	} else {
		d.event.Emit("download:completed", fmt.Sprintf("Downloaded %d files", downloadedMediaCount))
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}

func (d *Downloader) DownloadMedia(ctx context.Context, src string, dir string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, src, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	req.Header = http.Header{
		"User-Agent": {
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148 Safari/537.36",
		},
		"Accept":          {"image/avif,image/webp,image/apng,image/*,*/*;q=0.8"},
		"Accept-Language": {"en-US,en;q=0.9"},
		"Referer":         {deriveReferer(src)},
		"Connection":      {"keep-alive"},
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// ok
	case http.StatusForbidden:
		return fmt.Errorf("forbidden (headers/cookies/hotlink): %d, %s", resp.StatusCode, src)
	case http.StatusNotFound:
		return fmt.Errorf("media not found: %d, %s", resp.StatusCode, src)
	default:
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := d.SaveMedia(resp, src, dir); err != nil {
		return fmt.Errorf("save media: %w", err)
	}

	return nil
}

func (d *Downloader) SaveMedia(resp *http.Response, src string, dir string) error {
	fileName := filepath.Base(strings.Split(src, "?")[0])
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

func deriveReferer(src string) string {
	u, err := url.Parse(src)
	if err != nil {
		return "https://example.com/"
	}

	return u.Scheme + "://" + u.Host + "/"
}
