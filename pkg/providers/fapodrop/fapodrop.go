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

func (p *Provider) GetPhotos() ([]string, error) {
	if p.MinPhotoID > p.MaxPhotoID {
		return []string{}, fmt.Errorf("min photo ID (%d) is greater than max photo ID (%d)", p.MinPhotoID, p.MaxPhotoID)
	}

	recentPhotoID, err := p.GetRecentPhotoID()
	if err != nil {
		return []string{}, err
	}

	if p.MaxPhotoID > recentPhotoID {
		p.MaxPhotoID = recentPhotoID
	}

	photos := make([]string, p.MaxPhotoID-p.MinPhotoID+1)

	for i := p.MinPhotoID; i <= p.MaxPhotoID; i++ {
		photos[i-p.MinPhotoID] = strconv.Itoa(i)
	}

	return photos, nil
}

func (p *Provider) GetPhotoURL(photoID string) (string, error) {
	intPhotoID, err := strconv.Atoi(photoID)
	if err != nil {
		return "", err
	}

	urlWithoutID, err := buildURL(p.BaseURL, p.Username)
	if err != nil {
		return "", err
	}

	paddedID := fmt.Sprintf("%04d", intPhotoID)
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
