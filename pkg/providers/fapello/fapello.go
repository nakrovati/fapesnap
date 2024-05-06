package fapello

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Provider struct {
	ProviderName string
	BaseURL      string
}

func (p *Provider) InitProvider() {
	p.ProviderName = "fapello"
	p.BaseURL = "https://fapello.com"
}

func (p *Provider) GetProviderName() string {
	return p.ProviderName
}

func (p *Provider) GetPhotoURL(photoID int, userName string) (string, error) {
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

func (p *Provider) GetFileName(url string) string {
	parts := strings.Split(url, "/")

	return parts[len(parts)-1]
}

func (p *Provider) GetRecentPhotoID(name string) (int, error) {
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

func buildURL(baseURL string, userName string, recentID int) (string, error) {
	firstSymbol := string(userName[0])
	secondSymbol := string(userName[1])
	photoCountGroup := roundUp(recentID)

	photoURL, err := url.JoinPath(baseURL, "content", firstSymbol, secondSymbol, userName, strconv.Itoa(photoCountGroup))
	if err != nil {
		return "", err
	}

	return photoURL, nil
}

func parsePhotoID(url string) (int, error) {
	re := regexp.MustCompile(`\/(\d+)/$`)

	matches := re.FindStringSubmatch(url)

	if len(matches) < 2 {
		return 0, fmt.Errorf("invalid url: %s", url)
	}

	photoID, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	return photoID, nil
}

func roundUp(num int) int {
	if num < 1000 {
		return 1000
	}

	return ((num + 999) / 1000) * 1000
}
