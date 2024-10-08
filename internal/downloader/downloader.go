package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

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
) error {
	downloadDir, err := utils.GetDownloadDirectory(providerName, collectionName)
	if err != nil {
		return fmt.Errorf("failed to get download directory: %w", err)
	}

	runtime.EventsEmit(ctx, "download-start")

	downloadedPhotosCount := 0

	defer func() {
		runtime.EventsEmit(ctx, "download-complete", fmt.Sprintf("Downloaded %d photos", downloadedPhotosCount))
	}()

	for i := len(urls) - 1; i >= 0; i-- {
		select {
		case <-ctx.Done():
			runtime.EventsEmit(ctx,
				"download-canceled",
				fmt.Sprintf("Canceled after downloading %d photos", downloadedPhotosCount),
			)

			return ctx.Err()
		default:
			err = d.DownloadPhoto(ctx, urls[i], downloadDir)
			if err != nil {
				fmt.Printf("Failed to download photo from %s: %v\n", urls[i], err)

				continue
			}

			downloadedPhotosCount++
		}
	}

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
