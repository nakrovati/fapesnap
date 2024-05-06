package downloader

import (
	"fapesnap/pkg/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type PhotosProvider interface {
	InitProvider()
	GetProviderName() string
	GetFileName(src string) string
	GetPhotoURL(photoID int, userName string) (string, error)
	GetRecentPhotoID(userName string) (int, error)
}

type Downloader struct {
	PhotosProvider
}

func (d *Downloader) DownloadPhotos(userName string, min int, max int) error {
	if userName == "" {
		return fmt.Errorf("username cannot be empty")
	}

	err := utils.ValidateMinMax(min, max)
	if err != nil {
		return fmt.Errorf("invalid min/max: %w", err)
	}

	recentPhotoID := max
	if max == 100000 {
		if recentPhotoID, err = d.PhotosProvider.GetRecentPhotoID(userName); err != nil {
			return fmt.Errorf("failed to get recent photo ID: %w", err)
		}
	}

	downloadDir, err := utils.GetDownloadDirectory(d.PhotosProvider.GetProviderName(), userName)
	if err != nil {
		return fmt.Errorf("failed to get download directory: %w", err)
	}

	for i := recentPhotoID; i >= min; i-- {
		photoURL, err := d.PhotosProvider.GetPhotoURL(i, userName)
		if err != nil {
			return err
		}

		err = d.DownloadPhoto(photoURL, downloadDir)
		if err != nil {
			log.Println("error while downloading photo:", err)
			continue
		}

		fmt.Println("downloaded:", photoURL)
	}

	return nil
}

func (d Downloader) DownloadPhoto(src string, dir string) error {
	resp, err := http.Get(src)
	if err != nil {
		return fmt.Errorf("failed to download photo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("photo %s not found", src)
	}

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
