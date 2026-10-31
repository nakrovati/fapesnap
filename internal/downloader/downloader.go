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

func NewDownloader(logger *slog.Logger) *Downloader {
	if logger == nil {
		logger = slog.Default()
	}

	return &Downloader{
		logger: logger,
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

	emitEvent(d.logger, DownloadBatchStarted, struct{}{})

	jobs := make(chan providers.Media)

	var (
		wg                   sync.WaitGroup
		downloadedMediaCount int64
		failedMediaCount     int64
	)

	worker := func() {
		for media := range jobs {
			err := d.DownloadMedia(ctx, media.PageURL, media.URL, downloadDir)
			if err != nil {
				atomic.AddInt64(&failedMediaCount, 1)

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
