package fapodrop

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

func (p *Provider) GetProviderName() string {
	return p.ProviderName
}

func (p *Provider) InitProvider() {
	p.ProviderName = "fapodrop"
	p.BaseURL = "https://fapodrop.com"
}

func (p *Provider) GetCollectionName() string {
	return p.Username
}

func (p *Provider) GetMinMaxPhotoID() (int, int) {
	return p.MinPhotoID, p.MaxPhotoID
}

func (p *Provider) GetPhotoURL(photoID int) (string, error) {
	urlWithoutID, err := buildURL(p.BaseURL, p.Username)
	if err != nil {
		return "", err
	}

	paddedID := fmt.Sprintf("%04d", photoID)
	photoName := fmt.Sprintf("%s_%s.jpeg", p.Username, paddedID)

	url, err := url.JoinPath(urlWithoutID, photoName)
	if err != nil {
		return "", err
	}

	return url, nil
}

func buildURL(baseURL string, name string) (string, error) {
	firstSymbol := name[0]
	secondSymbol := name[1]

	photoURL, err := url.JoinPath(baseURL, "images", string(firstSymbol), string(secondSymbol), name, "1", "photo")
	if err != nil {
		return "", err
	}

	return photoURL, nil
}

func (p *Provider) GetRecentPhotoID() (int, error) {
	c := colly.NewCollector()

	userSrc, err := url.JoinPath(p.BaseURL, p.Username)
	if err != nil {
		return 0, err
	}

	isFound := false
	recentPhotoID := 0

	c.OnHTML(fmt.Sprintf(".row .one-pack a[href^='/%s']", p.Username), func(e *colly.HTMLElement) {
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

	return recentPhotoID, nil
}

func (p *Provider) GetFileName(url string) string {
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

	photoID, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}

	return photoID, err
}
