package downloader

import (
	"fapesnap/pkg/utils"
	"fmt"
	"io"
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
		return err
	}

	recentPhotoID := max
	if max == 100000 {
		if recentPhotoID, err = d.PhotosProvider.GetRecentPhotoID(userName); err != nil {
			return err
		}
	}

	downloadDir, err := utils.GetDownloadDirectory(d.PhotosProvider.GetProviderName(), userName)
	if err != nil {
		fmt.Println("Error while getting download directory:", err)
		return err
	}

	for i := recentPhotoID; i >= min; i-- {
		photoURL, err := d.PhotosProvider.GetPhotoURL(i, userName)
		if err != nil {
			return err
		}

		d.DownloadPhoto(photoURL, downloadDir)
		println("Downloaded:", photoURL)
	}
	return nil
}

func (d Downloader) DownloadPhoto(src string, dir string) {
	resp, err := http.Get(src)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	fileName := filepath.Join(dir, d.PhotosProvider.GetFileName(src))
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating file:", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error while copying file:", err)
		return
	}
}
