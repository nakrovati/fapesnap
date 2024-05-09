package downloader

import (
	"context"
	"errors"
	"fapesnap/internal/pkg/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const MaxPhotoID = 100000

var (
	ErrPhotoNotFound = errors.New("photo not found")
	ErrUsernameEmpty = errors.New("username cannot be empty")
)

type PhotosProvider interface {
	InitProvider()
	GetCollectionName() string
	GetProviderName() string
	GetFileName(src string) string
	GetPhotoURL(photoID string) (string, error)
	GetPhotos() ([]string, error)
	GetRecentPhotoID() (string, error)
}

type Downloader struct {
	PhotosProvider PhotosProvider
}

func (d *Downloader) DownloadPhotos() error {
	downloadDir, err := utils.GetDownloadDirectory(d.PhotosProvider.GetProviderName(), d.PhotosProvider.GetCollectionName())
	if err != nil {
		return fmt.Errorf("failed to get download directory: %w", err)
	}

	photos, err := d.PhotosProvider.GetPhotos()
	if err != nil {
		return fmt.Errorf("failed to get photos: %w", err)
	}

	for i := len(photos) - 1; i >= 0; i-- {
		photoID := photos[i]

		photoURL, err := d.PhotosProvider.GetPhotoURL(photoID)
		if err != nil {
			return fmt.Errorf("failed to get photo URL: %w", err)
		}

		err = d.DownloadPhoto(photoURL, downloadDir)
		if err != nil {
			log.Println("error while downloading photo:", err)

			continue
		}

		fmt.Printf("downloaded: %s\n", photoURL)
	}

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
	fileName := filepath.Join(dir, d.PhotosProvider.GetFileName(src))

	file, err := os.Create(fileName)
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
