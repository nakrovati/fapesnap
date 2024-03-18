package fapello

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
	providerName = "fapello"
	baseURL      = "https://fapello.com"
)

type FapelloProvider struct{}

func (p *FapelloProvider) DownloadPhotos(userName string, min int, max int) error {
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

	urlWithoutID, err := buildURL(userName, recentPhotoID)
	if err != nil {
		return err
	}

	for i := recentPhotoID; i >= 1; i-- {
		photoID := strconv.Itoa(i)
		photoName := fmt.Sprintf("%s_%s.jpg", userName, photoID)

		photoSrc, err := url.JoinPath(urlWithoutID, photoName)
		if err != nil {
			return err
		}

		println(photoSrc)
		downloadPhoto(photoSrc, downloadDir)

		println("Downloaded:", photoName)
	}

	return nil
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

func getRecentPhotoID(userName string) (int, error) {
	c := colly.NewCollector()

	recentPhotoSrc, err := url.JoinPath(baseURL, userName)
	if err != nil {
		return 0, err
	}

	isFound := false
	recentPhotoID := 0

	c.OnHTML("#content div a", func(e *colly.HTMLElement) {
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
		return 0, fmt.Errorf("user %s not found", userName)
	}

	return recentPhotoID, nil
}

func buildURL(userName string, recentId int) (string, error) {
	firstSymbol := userName[0]
	secondSymbol := userName[1]
	photoCountGroup := roundUp(recentId)

	url, err := url.JoinPath(baseURL, "content", string(firstSymbol), string(secondSymbol), userName, strconv.Itoa(photoCountGroup))
	if err != nil {
		return "", err
	}
	return url, nil
}

func parsePhotoID(url string) (int, error) {
	re := regexp.MustCompile(`\/(\d+)/$`)

	matches := re.FindStringSubmatch(url)

	if len(matches) < 2 {
		return 0, fmt.Errorf("invalid url: %s", url)
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	return id, nil

}

func roundUp(num int) int {
	if num < 1000 {
		return 1000
	}
	return ((num + 999) / 1000) * 1000
}
