package providers

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/nakrovati/fapesnap/internal/pkg/utils"
)

type FapelloProvider struct {
	ProviderName string
	BaseURL      string
}

func NewFapelloProvider() *FapelloProvider {
	return &FapelloProvider{
		ProviderName: "fapello",
		BaseURL:      "https://fapello.com",
	}
}

func (p *FapelloProvider) FetchMediaItems(collectionSlug string) ([]Media, error) {
	mediaItems := make([]Media, 0)

	for page := 1; ; page++ {
		pageURL := p.pageURL(collectionSlug, page)

		found, err := p.fetchProfilePage(pageURL, collectionSlug, &mediaItems)
		if err != nil {
			if errors.Is(err, ErrVisitNotFound) {
				return nil, fmt.Errorf("Profile not found: %s", collectionSlug)
			}

			return nil, err
		}

		if found == 0 {
			break
		}
	}

	if len(mediaItems) == 0 {
		return []Media{}, ErrNoMediaFound
	}

	return mediaItems, nil
}

func (p *FapelloProvider) GetCollectionFromURL(inputURL string) (string, error) {
	_, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	if !strings.Contains(inputURL, p.BaseURL) {
		return "", fmt.Errorf("%w: %s", ErrInvalidDomain, inputURL)
	}

	inputURL = strings.TrimSuffix(inputURL, "/")
	parts := strings.Split(inputURL, "/")

	if len(parts) < 4 || parts[len(parts)-1] == "" {
		return "", errors.New("can't get collection from url")
	}

	return parts[len(parts)-1], nil
}

func (p *FapelloProvider) getMedia(mediaID string, username string) (Media, error) {
	intMediaID, err := strconv.Atoi(mediaID)
	if err != nil {
		return Media{}, err
	}

	paddedID := fmt.Sprintf("%04d", intMediaID)
	mediaName := fmt.Sprintf("%s_%v.jpg", username, paddedID)
	mediaThumbnailName := fmt.Sprintf("%s_%v_300px.jpg", username, paddedID)

	urlWithoutID, err := p.buildURL(p.BaseURL, username, intMediaID)
	if err != nil {
		return Media{}, err
	}

	mediaURL, err := url.JoinPath(urlWithoutID, mediaName)
	if err != nil {
		return Media{}, err
	}

	thumbnailURL, err := url.JoinPath(urlWithoutID, mediaThumbnailName)
	if err != nil {
		return Media{}, err
	}

	media := Media{
		Type:         MediaTypeImage,
		URL:          mediaURL,
		ThumbnailURL: thumbnailURL,
	}

	return media, nil
}

func (p *FapelloProvider) getVideoURL(mediaPageURL string) (string, error) {
	c := colly.NewCollector()

	var (
		videoURL string
		visitErr error
	)

	c.OnHTML("video source[src]", func(e *colly.HTMLElement) {
		src := e.Attr("src")
		if src == "" {
			return
		}

		videoURL = src
	})

	c.OnError(func(c *colly.Response, err error) {
		visitErr = normalizeCollyError(c, err)
	})

	err := c.Visit(mediaPageURL)
	if visitErr != nil {
		return "", visitErr
	}

	if err != nil {
		return "", err
	}

	return videoURL, nil
}

func (p *FapelloProvider) buildURL(baseURL string, username string, recentID int) (string, error) {
	firstSymbol := string(username[0])
	secondSymbol := string(username[1])
	countGroup := strconv.Itoa(utils.RoundUp(recentID))

	mediaURL, err := url.JoinPath(baseURL, "content", firstSymbol, secondSymbol, username, countGroup)
	if err != nil {
		return "", fmt.Errorf("build url failed: %w", err)
	}

	return mediaURL, nil
}

func (p *FapelloProvider) parseMediaID(url string) (int, error) {
	re := regexp.MustCompile(`\/(\d+)/$`)

	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return 0, fmt.Errorf("invalid url: %s", url)
	}

	mediaID, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	return mediaID, nil
}

func (p *FapelloProvider) pageURL(slug string, page int) string {
	if page == 1 {
		return fmt.Sprintf("%s/%s/", p.BaseURL, slug)
	}

	return fmt.Sprintf("%s/ajax/model/%s/page-%d/", p.BaseURL, slug, page)
}

func (p *FapelloProvider) fetchProfilePage(targetURL string, collectionSlug string, mediaItems *[]Media) (int, error) {
	c := colly.NewCollector()

	found := 0         // Page load requests (/ajax/model/...) return a 200 status code
	var visitErr error //

	c.OnHTML(
		fmt.Sprintf("a[href*='/%s/']", collectionSlug),
		func(e *colly.HTMLElement) {
			media, ok := p.parseCard(e, collectionSlug)
			if !ok {
				return
			}

			*mediaItems = append(*mediaItems, media)

			found++
		},
	)

	c.OnError(func(c *colly.Response, err error) {
		visitErr = normalizeCollyError(c, err)
	})

	err := c.Visit(targetURL)
	if visitErr != nil {
		return 0, visitErr
	}

	if err != nil {
		return 0, err
	}

	return found, nil
}

func (p *FapelloProvider) isVideoCard(e *colly.HTMLElement) bool {
	return strings.Contains(e.Text, "icon-play.svg") ||
		e.ChildAttr("img[src*='icon-play.svg']", "src") != ""
}

func (p *FapelloProvider) parseCard(e *colly.HTMLElement, slug string) (Media, bool) {
	href := e.Attr("href")
	if href == "" {
		return Media{}, false
	}

	img := e.ChildAttr("img", "src")
	if img == "" {
		return Media{}, false
	}

	id, err := p.parseMediaID(href)
	if err != nil {
		return Media{}, false
	}

	media, err := p.getMedia(strconv.Itoa(id), slug)
	if err != nil {
		return Media{}, false
	}

	media.ThumbnailURL = img

	if p.isVideoCard(e) {
		videoURL, err := p.getVideoURL(href)
		if err != nil {
			return Media{}, false
		}

		media.Type = MediaTypeVideo
		media.URL = videoURL
	} else {
		media.Type = MediaTypeImage
	}

	return media, true
}
