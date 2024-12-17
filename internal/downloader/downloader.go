package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/nakrovati/fapesnap/internal/pkg/utils"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const MaxPhotoID = 100000

var (
	ErrPhotoNotFound = errors.New("photo not found")
	ErrUsernameEmpty = errors.New("username cannot be empty")
)

type Downloader struct{}

func (d *Downloader) DownloadPhotos(
	ctx context.Context,
	urls []string,
	providerName string,
	collectionName string,
	maxParallelDownloads int,
) error {
	semaphore := make(chan struct{}, maxParallelDownloads)

	var wg sync.WaitGroup

	counterChan := make(chan int)

	var downloadedPhotosCount int

	go func() {
		for count := range counterChan {
			downloadedPhotosCount += count
		}
	}()

	downloadDir, err := utils.GetDownloadDirectory(providerName, collectionName)
	if err != nil {
		return fmt.Errorf("failed to get download directory: %w", err)
	}

	runtime.EventsEmit(ctx, "download-start")

	defer func() {
		close(counterChan)
		runtime.EventsEmit(ctx, "download-complete", fmt.Sprintf("Downloaded %d photos", downloadedPhotosCount))
	}()

	for _, url := range urls {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				fmt.Printf("Download cancelled for %s\n", url)
			}

			err := d.DownloadPhoto(ctx, url, downloadDir)
			if err != nil {
				fmt.Printf("Failed to download photo from %s: %v\n", url, err)
			} else {
				fmt.Printf("Downloaded %s\n", url)
				counterChan <- 1
			}
		}(url)
	}

	wg.Wait()

	return nil
}

func (d Downloader) DownloadPhoto(ctx context.Context, src string, dir string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, src, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download photo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("photo %s not found %w", src, ErrPhotoNotFound)
	}

	err = d.SavePhoto(resp, src, dir)
	if err != nil {
		return err
	}

	return nil
}

func (d Downloader) SavePhoto(resp *http.Response, src string, dir string) error {
	fileName := filepath.Base(src)
	filePath := filepath.Join(dir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save photo: %w", err)
	}

	return nil
}
