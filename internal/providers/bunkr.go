package providers

import (
	"errors"
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
	mediaItems := make([]Media, 0)

	err := p.fetchAlbumPage(collectionSlug, &mediaItems)
	if err != nil {
		return []Media{}, err
	}

	if len(mediaItems) == 0 {
		return []Media{}, errors.New("no media found")
	}

	return mediaItems, nil
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

func (p *BunkrProvider) getMediaURL(href string) (string, error) {
	c := colly.NewCollector()

	var mediaURL string

	c.OnHTML("main.cont", func(e *colly.HTMLElement) {
		mediaURL = e.ChildAttr(
			".cont img.w-full.h-full.absolute[src]",
			"src",
		)
		if mediaURL != "" {
			return
		}

		downloadURL := e.ChildAttr(
			"a[href*=\"get.bunkrr.su/file/\"]",
			"href",
		)

		if downloadURL != "" {
			mediaURL = downloadURL + "#"
		}
	})

	err := c.Visit(href)
	if err != nil {
		return "", err
	}

	return mediaURL, nil
}

func (p *BunkrProvider) fetchAlbumPage(collectionSlug string, mediaItems *[]Media) error {
	albumURL, err := url.JoinPath(p.BaseURL, "a", collectionSlug)
	if err != nil {
		return err
	}

	c := colly.NewCollector()

	c.OnHTML(".theItem", func(e *colly.HTMLElement) {
		item := p.parseItem(e)

		mediaURL, err := p.getMediaURL(item.URL)
		if err != nil {
			return
		}

		item.URL = mediaURL

		*mediaItems = append(*mediaItems, item)
	})

	err = c.Visit(albumURL)
	if err != nil {
		return err
	}

	if len(*mediaItems) == 0 {
		return errors.New("no media resources were found")
	}

	return nil
}

func (p *BunkrProvider) parseItem(e *colly.HTMLElement) Media {
	mediaType := p.checkMediaType(e)

	href := e.ChildAttr("a[aria-label='download']", "href")

	var thumbnailURL string
	if mediaType == MediaTypeImage || mediaType == MediaTypeVideo {
		thumbnailURL = e.ChildAttr("img.grid-images_box-img", "src")
	}

	itemURL, err := url.JoinPath(p.BaseURL, href)
	if err != nil {
		return Media{}
	}

	return Media{
		URL:          itemURL,
		Name:         e.ChildText(".theName"),
		Size:         e.ChildText(".theSize"),
		ThumbnailURL: thumbnailURL,
		Type:         mediaType,
	}
}

func (p *BunkrProvider) checkMediaType(e *colly.HTMLElement) MediaType {
	class := e.ChildAttr(`span[class*="type-"]`, "class")

	switch {
	case strings.Contains(class, "type-Image"):
		return MediaTypeImage
	case strings.Contains(class, "type-Video"):
		return MediaTypeVideo
	case strings.Contains(class, "type-File"):
		return MediaTypeFile
	default:
		return MediaTypeUnknown
	}
}
