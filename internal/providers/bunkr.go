package providers

import (
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/gocolly/colly/v2"
)

type BunkrProvider struct {
	ProviderName string
	BaseURL      string
}

func (p *BunkrProvider) InitProvider() {
	p.ProviderName = "bunkr"
	p.BaseURL = "https://bunkrrr.org"
}

func (p BunkrProvider) FetchPhotoURLs(collection string) ([]string, error) {
	photosURLsFromHref, err := p.GetPhotoURLs(collection)
	if err != nil {
		return []string{}, err
	}

	photos := make([]string, 0, len(photosURLsFromHref))

	for _, photoURLFromHref := range photosURLsFromHref {
		photoURL, err := p.GetPhotoURL(photoURLFromHref)
		if err != nil {
			return []string{}, err
		}

		photos = append(photos, photoURL)
	}

	slices.Reverse(photos)

	return photos, nil
}

func (p BunkrProvider) GetPhotoURLs(albumID string) ([]string, error) {
	albumURL, err := url.JoinPath(p.BaseURL, "a", albumID)
	if err != nil {
		return []string{}, err
	}

	photosURLs := make([]string, 0)

	c := colly.NewCollector()

	c.OnHTML("a.grid-images_box-link", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		photosURLs = append(photosURLs, href)
	})

	c.Visit(albumURL)

	if len(photosURLs) == 0 {
		return []string{}, fmt.Errorf("album not found")
	}

	return photosURLs, nil
}

func (p BunkrProvider) GetCollectionFromURL(inputURL string) (string, error) {
	_, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	if !strings.Contains(inputURL, p.BaseURL) && !strings.Contains(inputURL, "bunkr") {
		return "", errors.New("Unvalid domain")
	}

	inputURL = strings.TrimSuffix(inputURL, "/")
	parts := strings.Split(inputURL, "/")

	if len(parts) < 5 || parts[len(parts)-1] == "" {
		return "", errors.New("Can't get collection from url")
	}

	return parts[len(parts)-1], nil
}

func (p BunkrProvider) GetPhotoURL(photoURL string) (string, error) {
	c := colly.NewCollector()
	c.OnHTML(".lightgallery img", func(e *colly.HTMLElement) {
		if src := e.Attr("src"); src != "" {
			photoURL = src
		}
	})

	c.Visit(photoURL)

	if photoURL == "" {
		return "", fmt.Errorf("photo not found")
	}

	return photoURL, nil
}

func (p BunkrProvider) GetPhotoID(src string) string {
	u, err := url.Parse(src)
	if err != nil {
		return ""
	}

	pathParts := strings.Split(u.Path, "/")

	id := pathParts[len(pathParts)-1]

	return id
}
