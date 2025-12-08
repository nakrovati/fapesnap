package providers

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
)

type BunkrProvider struct {
	ProviderName string
	BaseURL      string
}

func NewBunkrProvider() *BunkrProvider {
	return &BunkrProvider{
		ProviderName: "bunkr",
		BaseURL:      "https://bunkr.cr",
	}
}

func (p *BunkrProvider) FetchPhotoURLs(collection string) ([]Photo, error) {
	items, err := p.GetPhotos(collection)
	if err != nil {
		return []Photo{}, err
	}

	photos := make([]Photo, 0, len(items))

	for _, item := range items {
		photoURL, err := p.GetPhotoURL(item.Href)
		if err != nil {
			fmt.Printf("Failed to get photo: %v\n", err)

			continue
		}

		photo := Photo{
			URL:          photoURL,
			ThumbnailURL: item.ThumbnailURL,
		}

		photos = append(photos, photo)
	}

	if len(photos) == 0 {
		return []Photo{}, errors.New("no photos found")
	}

	return photos, nil
}

type BunkrItem struct {
	Href         string
	ThumbnailURL string
}

func (p *BunkrProvider) GetPhotos(albumID string) ([]BunkrItem, error) {
	albumURL, err := url.JoinPath(p.BaseURL, "a", albumID)
	if err != nil {
		return []BunkrItem{}, err
	}

	items := make([]BunkrItem, 0)

	c := colly.NewCollector()

	c.OnHTML(".theItem", func(e *colly.HTMLElement) {
		photoPageURL := e.ChildAttr("a[aria-label='download']", "href")
		thumbnailURL := e.ChildAttr("img.grid-images_box-img", "src")

		href, err := url.JoinPath(p.BaseURL, photoPageURL)
		if err != nil {
			fmt.Println(err)

			return
		}

		item := BunkrItem{
			Href:         href,
			ThumbnailURL: thumbnailURL,
		}

		items = append(items, item)
	})

	err = c.Visit(albumURL)
	if err != nil {
		return []BunkrItem{}, err
	}

	if len(items) == 0 {
		return []BunkrItem{}, errors.New("album not found")
	}

	return items, nil
}

func (p *BunkrProvider) GetCollectionFromURL(inputURL string) (string, error) {
	_, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	if !strings.Contains(inputURL, p.BaseURL) && !strings.Contains(inputURL, "bunkr") {
		return "", errors.New("unvalid domain")
	}

	inputURL = strings.TrimSuffix(inputURL, "/")
	parts := strings.Split(inputURL, "/")

	if len(parts) < 5 || parts[len(parts)-1] == "" {
		return "", errors.New("can't get collection from url")
	}

	return parts[len(parts)-1], nil
}

func (p *BunkrProvider) GetPhotoURL(photoURL string) (string, error) {
	c := colly.NewCollector()

	c.OnHTML("main.cont", func(e *colly.HTMLElement) {
		photoURL = e.ChildAttr("img.w-full.h-full.absolute", "src")
	})

	err := c.Visit(photoURL)
	if err != nil {
		return "", err
	}

	return photoURL, nil
}

func (p *BunkrProvider) GetPhotoID(src string) string {
	u, err := url.Parse(src)
	if err != nil {
		return ""
	}

	pathParts := strings.Split(u.Path, "/")
	id := pathParts[len(pathParts)-1]

	return id
}
