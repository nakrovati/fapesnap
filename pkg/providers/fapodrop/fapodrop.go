package fapodrop

import (
	"fapesnap/pkg/utils"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var (
	providerName = "fapodrop"
	baseURL      = "https://www.fapodrop.com"
)

type FapodropProvider struct{}

func (p *FapodropProvider) DownloadPhotos(userName string, min int, max int) error {
	downloadDir, err := utils.GetDownloadDirectory(providerName, userName)
	if err != nil {
		fmt.Println("Error while getting download directory:", err)
		return err
	}

	recentPhotoID := max
	if max == 100000 {
		if recentPhotoID, err = getRecentPhotoID(userName); err != nil {
			return err
		}
	}

	urlWithoutID, err := buildURL(userName)
	if err != nil {
		return err
	}

	for i := recentPhotoID; i >= min; i-- {
		paddedID := fmt.Sprintf("%04d", i)
		photoName := fmt.Sprintf("%s_%s.jpeg", userName, paddedID)

		photoSrc, err := url.JoinPath(urlWithoutID, photoName)
		if err != nil {
			return err
		}

		downloadPhoto(photoSrc, downloadDir)

		println("Downloaded:", photoName)
	}
	return nil
}

func buildURL(name string) (string, error) {
	firstSymbol := name[0]
	secondSymbol := name[1]

	return url.JoinPath(baseURL, "images", string(firstSymbol), string(secondSymbol), name, "1", "photo")
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

			photoID, err := parsePhotoID(src)
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

func getFileName(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

func parsePhotoID(url string) (int, error) {
	re := regexp.MustCompile(`\/\d{4}$`)

	match := re.FindString(url)
	if match == "" {
		return 0, fmt.Errorf("invalid url: %s", url)
	}

	numStr := match[1:] // Take out the first "/"
	return strconv.Atoi(numStr)
}
