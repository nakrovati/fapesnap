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

type Downloader struct {
	ctx context.Context
}

func (d *Downloader) SetContext(ctx context.Context) {
	d.ctx = ctx
}

func (d *Downloader) DownloadPhotos(urls []string, providerName string, collectionName string) error {
	downloadDir, err := utils.GetDownloadDirectory(providerName, collectionName)
	if err != nil {
		return fmt.Errorf("failed to get download directory: %w", err)
	}

	runtime.EventsEmit(d.ctx, "download-start")

	downloadedPhotosCount := 0

	for i := len(urls) - 1; i >= 0; i-- {
		err = d.DownloadPhoto(urls[i], downloadDir)
		if err != nil {
			continue
		}
		downloadedPhotosCount += 1
	}

	runtime.EventsEmit(d.ctx, "download-complete", fmt.Sprintf("Downloaded %d photos", downloadedPhotosCount))

	return nil
}

func (d Downloader) DownloadPhoto(src string, dir string) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, src, nil)
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
