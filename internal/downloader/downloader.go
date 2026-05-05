package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/nakrovati/fapesnap/internal/config"
	"github.com/nakrovati/fapesnap/internal/pkg/utils"
	"github.com/nakrovati/fapesnap/internal/providers"
	"github.com/wailsapp/wails/v3/pkg/application"
)

var (
	ErrPhotoNotFound = errors.New("photo not found")
	ErrUsernameEmpty = errors.New("username cannot be empty")
)

type Downloader struct {
	httpClient *http.Client
}

func NewDownloader() *Downloader {
	return &Downloader{
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (d *Downloader) DownloadPhotos(
	app *application.App,
	ctx context.Context,
	photos []providers.Photo,
	baseDownloadDir config.DownloadDir,
	providerName string,
	collectionSlug string,
	maxParallelDownloads int,
) error {
	downloadDir, err := utils.GetCollectionDownloadDir(baseDownloadDir, providerName, collectionSlug)
	if err != nil {
		return fmt.Errorf("failed to get download directory: %w", err)
	}

	app.Event.Emit("download-start")

	jobs := make(chan providers.Photo)
	counterChan := make(chan int)

	var wg sync.WaitGroup

	var downloadedPhotosCount int

	go func() {
		for count := range counterChan {
			downloadedPhotosCount += count
		}
	}()

	for range maxParallelDownloads {
		wg.Go(func() {
			for photo := range jobs {
				if ctx.Err() != nil {
					fmt.Printf("Download cancelled for %s\n", photo.URL)

					continue
				}

				err := d.DownloadPhoto(ctx, photo.URL, downloadDir)
				if err != nil {
					fmt.Printf("Failed to download photo: %v\n", err)

					continue
				}

				fmt.Printf("Downloaded %s\n", photo.URL)

				counterChan <- 1
			}
		})
	}

	for _, photo := range photos {
		if ctx.Err() != nil {
			break
		}

		jobs <- photo
	}

	close(jobs)
	wg.Wait()

	close(counterChan)
	app.Event.Emit("download-complete",
		fmt.Sprintf("Downloaded %d photos", downloadedPhotosCount),
	)

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}

func (d *Downloader) DownloadPhoto(ctx context.Context, src string, dir string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, src, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
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
		return fmt.Errorf("failed to download photo: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// ok
	case http.StatusForbidden:
		return fmt.Errorf("%d forbidden (headers/cookies/hotlink) for %s", resp.StatusCode, src)
	case http.StatusNotFound:
		return fmt.Errorf("%d photo %s not found: %w", resp.StatusCode, src, ErrPhotoNotFound)
	default:
		return fmt.Errorf("%d failed to download photo", resp.StatusCode)
	}

	err = d.SavePhoto(resp, src, dir)
	if err != nil {
		return err
	}

	return nil
}

func (d *Downloader) SavePhoto(resp *http.Response, src string, dir string) error {
	fileName := filepath.Base(strings.Split(src, "?")[0])

	filePath := filepath.Join(dir, fileName)

	if !strings.HasPrefix(filepath.Clean(filePath), filepath.Clean(dir)) {
		return fmt.Errorf("file path escapes target directory: %s", filePath)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			fmt.Printf("Failed to close file: %v\n", err)
		}
	}()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save photo: %w", err)
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
