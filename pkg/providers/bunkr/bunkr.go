package bunkr

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

type Provider struct {
	Album        string
	ProviderName string
	BaseURL      string
}

func (p *Provider) InitProvider() {
	p.ProviderName = "bunkr"
	p.BaseURL = "https://bunkr.si"
}

func (p *Provider) GetProviderName() string {
	return p.ProviderName
}

func (p *Provider) GetCollectionName() string {
	return p.Album
}

func (p *Provider) GetPhotoURL(photoID string) (string, error) {
	photoPageURL, err := url.JoinPath(p.BaseURL, "i", photoID)
	if err != nil {
		return "", err
	}

	c := colly.NewCollector()

	photoURL := ""

	c.OnHTML(".lightgallery img", func(e *colly.HTMLElement) {
		if src := e.Attr("src"); src != "" {
			photoURL = src
		}
	})

	c.Visit(photoPageURL)

	if photoURL == "" {
		return "", fmt.Errorf("photo not found")
	}

	return photoURL, nil
}

func (p *Provider) GetFileName(src string) string {
	parts := strings.Split(src, "/")

	return parts[len(parts)-1]
}

func (p *Provider) GetPhotos() ([]string, error) {
	albumURL, err := url.JoinPath(p.BaseURL, "a", p.Album)
	if err != nil {
		return []string{}, err
	}

	_, err = p.GetRecentPhotoID()
	if err != nil {
		return []string{}, err
	}

	photos := make([]string, 0)

	c := colly.NewCollector()

	c.OnHTML(".grid-images_box-link", func(e *colly.HTMLElement) {
		photos = append(photos, getPhotoID(e.Attr("href")))
	})

	c.Visit(albumURL)

	if len(photos) == 0 {
		return []string{}, fmt.Errorf("photos not found")
	}

	return photos, nil
}

func (p *Provider) GetRecentPhotoID() (string, error) {
	c := colly.NewCollector()

	albumURL, err := url.JoinPath(p.BaseURL, "a", p.Album)
	if err != nil {
		return "", err
	}

	recentPhotoID := ""
	isFound := false

	c.OnHTML(".grid-images_box-link", func(e *colly.HTMLElement) {
		if !isFound {
			href := e.Attr("href")
			recentPhotoID = getPhotoID(href)
		}
	})

	c.Visit(albumURL)

	if recentPhotoID == "" {
		return "", fmt.Errorf("album not found")
	}

	return recentPhotoID, nil
}

func getPhotoID(src string) string {
	u, err := url.Parse(src)
	if err != nil {
		return ""
	}

	pathParts := strings.Split(u.Path, "/")

	id := pathParts[len(pathParts)-1]

	return id

}
