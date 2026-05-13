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

func (p *BunkrProvider) FetchMediaItems(collectionSlug string) ([]Media, error) {
	items, err := p.GetMedias(collectionSlug)
	if err != nil {
		return []Media{}, err
	}

	mediaItems := make([]Media, 0, len(items))

	for _, item := range items {
		mediaURL, err := p.GetMediaURL(item.Href)
		if err != nil {
			fmt.Printf("Failed to get media: %v\n", err)

			continue
		}

		media := Media{
			URL:          mediaURL,
			ThumbnailURL: item.ThumbnailURL,
		}

		mediaItems = append(mediaItems, media)
	}

	if len(mediaItems) == 0 {
		return []Media{}, errors.New("no media found")
	}

	return mediaItems, nil
}

type BunkrItem struct {
	Href         string
	ThumbnailURL string
}

func (p *BunkrProvider) GetMedias(albumID string) ([]BunkrItem, error) {
	albumURL, err := url.JoinPath(p.BaseURL, "a", albumID)
	if err != nil {
		return []BunkrItem{}, err
	}

	items := make([]BunkrItem, 0)

	c := colly.NewCollector()

	c.OnHTML(".theItem", func(e *colly.HTMLElement) {
		mediaPageURL := e.ChildAttr("a[aria-label='download']", "href")
		thumbnailURL := e.ChildAttr("img.grid-images_box-img", "src")

		href, err := url.JoinPath(p.BaseURL, mediaPageURL)
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

func (p *BunkrProvider) GetMediaURL(mediaURL string) (string, error) {
	c := colly.NewCollector()

	c.OnHTML("main.cont", func(e *colly.HTMLElement) {
		mediaURL = e.ChildAttr("img.w-full.h-full.absolute", "src")
	})

	err := c.Visit(mediaURL)
	if err != nil {
		return "", err
	}

	return mediaURL, nil
}

func (p *BunkrProvider) GetMediaID(src string) string {
	u, err := url.Parse(src)
	if err != nil {
		return ""
	}

	pathParts := strings.Split(u.Path, "/")
	id := pathParts[len(pathParts)-1]

	return id
}
