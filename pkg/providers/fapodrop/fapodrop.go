package fapodrop

import (
	"fapesnap/pkg/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly"
)

var (
	providerName = "fapodrop"
	baseURL      = "https://www.fapodrop.com"
)

func buildURL(name string) (string, error) {
	firstSymbol := name[0]
	secondSymbol := name[1]

	url, err := url.JoinPath(baseURL, "images", string(firstSymbol), string(secondSymbol), name, "1", "photo")
	if err != nil {
		return "", err
	}
	return url, nil
}

func getRecentPhotoID(name string) (int, error) {
	c := colly.NewCollector()

	recentPhotoSrc, err := url.JoinPath(baseURL, name)
	if err != nil {
		return 0, err
	}

	isFound := false
	recentPhotoID := 0

	c.OnHTML(fmt.Sprintf(".row .one-pack a[href^='/%s']", name), func(e *colly.HTMLElement) {
		if !isFound {
			src := e.Attr("href")

			photoID, err := utils.ParsePhotoID(src)
			if err != nil {
				return
			}

			recentPhotoID = photoID
			isFound = true
		}
	})

	c.Visit(recentPhotoSrc)

	if !isFound {
		return 0, fmt.Errorf("user %s not found", name)
	}

	return recentPhotoID, nil
}

func downloadPhoto(src string, downloadDir string) {
	resp, err := http.Get(src)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	fileName := filepath.Join(downloadDir, getFileName(src))
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

func DownloadPhotos(name string) string {
	recentPhotoID, err := getRecentPhotoID(name)
	if err != nil {
		log.Fatal(err)
	}

	urlWithoutID, err := buildURL(name)
	if err != nil {
		return ""
	}

	downloadDir, err := getDownloadDirectory(name)
	if err != nil {
		fmt.Println("Error while getting download directory:", err)
		return ""
	}

	for i := recentPhotoID; i >= 1; i-- {
		paddedID := fmt.Sprintf("%04d", i)
		photoName := fmt.Sprintf("%s_%s.jpeg", name, paddedID)

		photoSrc, err := url.JoinPath(urlWithoutID, photoName)
		if err != nil {
			return ""
		}

		downloadPhoto(photoSrc, downloadDir)

		println("Downloaded:", photoName)
	}
	return ""
}

func getDownloadDirectory(name string) (string, error) {
	usr, err := user.Current()

	if err != nil {
		return "", err
	}
	downloadDir := filepath.Join(usr.HomeDir, "Downloads", "fapesnap")

	providerDir := filepath.Join(downloadDir, providerName, name)
	err = os.MkdirAll(providerDir, 0755)
	if err != nil {
		fmt.Println("Error while creating download directory:", err)
		return "", err
	}

	return providerDir, nil
}

func getFileName(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}
