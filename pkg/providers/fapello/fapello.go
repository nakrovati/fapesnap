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
	MaxPhotoID   int
	MinPhotoID   int
	ProviderName string
	BaseURL      string
	Username     string
}

func (p *Provider) InitProvider() {
	p.ProviderName = "fapello"
	p.BaseURL = "https://fapello.com"
}

func (p *Provider) GetMinMaxPhotoID() (int, int) {
	return p.MinPhotoID, p.MaxPhotoID
}

func (p *Provider) GetProviderName() string {
	return p.ProviderName
}

func (p *Provider) GetCollectionName() string {
	return p.Username
}

func (p *Provider) GetPhotoURL(photoID int) (string, error) {
	urlWithoutID, err := buildURL(p.BaseURL, p.Username, photoID)
	if err != nil {
		return "", err
	}

	paddedID := fmt.Sprintf("%04d", photoID)
	photoName := fmt.Sprintf("%s_%v.jpg", p.Username, paddedID)

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

func (p *Provider) GetRecentPhotoID() (int, error) {
	c := colly.NewCollector()

	userSrc, err := url.JoinPath(p.BaseURL, p.Username)
	if err != nil {
		return 0, err
	}

	isFound := false
	recentPhotoID := 0

	c.OnHTML(fmt.Sprintf("#content div a[href*='%s']", p.Username), func(e *colly.HTMLElement) {
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

	c.Visit(userSrc)

	if !isFound {
		return 0, fmt.Errorf("user not found")
	}

	return recentPhotoID, nil
}

func buildURL(baseURL string, username string, recentID int) (string, error) {
	firstSymbol := string(username[0])
	secondSymbol := string(username[1])
	photoCountGroup := strconv.Itoa(roundUp(recentID))

	photoURL, err := url.JoinPath(baseURL, "content", firstSymbol, secondSymbol, username, photoCountGroup)
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
