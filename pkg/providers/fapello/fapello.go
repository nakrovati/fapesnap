package fapello

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type FapelloProvider struct {
	ProviderName string
	BaseURL      string
}

func (p *FapelloProvider) InitProvider() {
	p.ProviderName = "fapello"
	p.BaseURL = "https://fapello.com"
}

func (p *FapelloProvider) GetProviderName() string {
	return p.ProviderName
}

func (p *FapelloProvider) GetPhotoURL(photoID int, userName string) (string, error) {
	urlWithoutID, err := buildURL(p.BaseURL, userName, photoID)
	if err != nil {
		return "", err
	}

	paddedID := fmt.Sprintf("%04d", photoID)
	photoName := fmt.Sprintf("%s_%v.jpg", userName, paddedID)

	url, err := url.JoinPath(urlWithoutID, photoName)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (p *FapelloProvider) GetFileName(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

func (p *FapelloProvider) GetRecentPhotoID(name string) (int, error) {
	c := colly.NewCollector()

	recentPhotoSrc, err := url.JoinPath(p.BaseURL, name)
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

	if err := c.Visit(recentPhotoSrc); err != nil {
		return 0, err
	}

	return recentPhotoID, nil
}

func buildURL(baseURL string, userName string, recentId int) (string, error) {
	firstSymbol := userName[0]
	secondSymbol := userName[1]
	photoCountGroup := roundUp(recentId)

	return url.JoinPath(baseURL, "content", string(firstSymbol), string(secondSymbol), userName, strconv.Itoa(photoCountGroup))
}

func parsePhotoID(url string) (int, error) {
	re := regexp.MustCompile(`\/(\d+)/$`)

	matches := re.FindStringSubmatch(url)

	if len(matches) < 2 {
		return 0, fmt.Errorf("invalid url: %s", url)
	}

	return strconv.Atoi(matches[1])
}

func roundUp(num int) int {
	if num < 1000 {
		return 1000
	}
	return ((num + 999) / 1000) * 1000
}
