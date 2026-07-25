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
	mediaItems := make([]Media, 0)

	err := p.fetchAlbumPage(collectionSlug, &mediaItems)
	if err != nil {
		if errors.Is(err, ErrVisitNotFound) {
			return nil, fmt.Errorf("profile not found: %s", collectionSlug)
		}

		return []Media{}, err
	}

	if len(mediaItems) == 0 {
		return []Media{}, ErrNoMediaFound
	}

	return mediaItems, nil
}

func (p *BunkrProvider) GetCollectionFromURL(inputURL string) (string, error) {
	_, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	if !strings.Contains(inputURL, p.BaseURL) && !strings.Contains(inputURL, "bunkr") {
		return "", fmt.Errorf("%w: %s", ErrInvalidDomain, inputURL)
	}

	inputURL = strings.TrimSuffix(inputURL, "/")
	parts := strings.Split(inputURL, "/")

	if len(parts) < 5 || parts[len(parts)-1] == "" {
		return "", errors.New("can't get collection from url")
	}

	return parts[len(parts)-1], nil
}

func (p *BunkrProvider) GetMediaURL(pageURL string, collectionSlug string) (string, error) {
	c := colly.NewCollector()

	var (
		mediaURL string
		visitErr error
	)

	c.OnHTML("main.cont", func(e *colly.HTMLElement) {
		mediaURL = e.ChildAttr(
			"img.w-full.h-full.absolute[src]",
			"src",
		)
	})

	c.OnError(func(c *colly.Response, err error) {
		visitErr = normalizeCollyError(c, err)
	})

	err := c.Visit(pageURL)

	if visitErr != nil {
		return "", visitErr
	}

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

	var visitErr error

	c.OnHTML(".theItem", func(e *colly.HTMLElement) {
		item := p.parseItem(e)

		*mediaItems = append(*mediaItems, item)
	})

	c.OnError(func(c *colly.Response, err error) {
		visitErr = normalizeCollyError(c, err)
	})

	err = c.Visit(albumURL)

	if visitErr != nil {
		return visitErr
	}

	if err != nil {
		return err
	}

	return nil
}

func (p *BunkrProvider) parseItem(e *colly.HTMLElement) Media {
	mediaType := p.checkMediaType(e)

	href := e.ChildAttr("a[aria-label='download']", "href")

	pageURL, err := url.JoinPath(p.BaseURL, href)
	if err != nil {
		return Media{}
	}

	var thumbnailURL string
	if mediaType == MediaTypeImage || mediaType == MediaTypeVideo {
		thumbnailURL = e.ChildAttr("img.grid-images_box-img", "src")
	}

	return Media{
		PageURL:      pageURL,
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
