package providers

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type FapodropProvider struct {
	ProviderName string
	BaseURL      string
}

func NewFapodropProvider() *FapodropProvider {
	return &FapodropProvider{
		ProviderName: "fapodrop",
		BaseURL:      "https://fapodrop.com",
	}
}

func (p *FapodropProvider) FetchMediaItems(collectionSlug string) ([]Media, error) {
	recentMediaID, err := p.GetRecentMediaID(collectionSlug)
	if err != nil {
		return []Media{}, err
	}

	minMediaID := 1
	maxMediaID := min(100000, recentMediaID)

	mediaItems := make([]Media, 0, maxMediaID)

	for i := maxMediaID; i >= minMediaID; i-- {
		media, err := p.GetMedia(strconv.Itoa(i), collectionSlug)
		if err != nil {
			fmt.Printf("Failed to get media: %v\n", err)

			continue
		}

		mediaItems = append(mediaItems, media)
	}

	if len(mediaItems) == 0 {
		return []Media{}, errors.New("no media found")
	}

	return mediaItems, nil
}

func (p *FapodropProvider) GetCollectionFromURL(inputURL string) (string, error) {
	_, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	if !strings.Contains(inputURL, p.BaseURL) {
		return "", errors.New("unvalid domain")
	}

	inputURL = strings.TrimSuffix(inputURL, "/")
	parts := strings.Split(inputURL, "/")

	if len(parts) < 4 || parts[len(parts)-1] == "" {
		return "", errors.New("can't get collection from url")
	}

	return parts[len(parts)-1], nil
}

func (p *FapodropProvider) GetMedia(mediaID string, username string) (Media, error) {
	intMediaID, err := strconv.Atoi(mediaID)
	if err != nil {
		return Media{}, err
	}

	paddedID := fmt.Sprintf("%04d", intMediaID)
	mediaName := fmt.Sprintf("%s_%s.jpeg", username, paddedID)
	mediaThumbnailName := fmt.Sprintf("%s_%s_thumbnail.jpeg", username, paddedID)

	urlWithoutID, err := p.buildURL(p.BaseURL, username)
	if err != nil {
		return Media{}, err
	}

	thumbnailURLWithoutID, err := p.buildThumbnailURL(p.BaseURL, username)
	if err != nil {
		return Media{}, err
	}

	mediaURL, err := url.JoinPath(urlWithoutID, mediaName)
	if err != nil {
		return Media{}, err
	}

	thumbnailURL, err := url.JoinPath(thumbnailURLWithoutID, mediaThumbnailName)
	if err != nil {
		return Media{}, err
	}

	media := Media{
		URL:          mediaURL,
		ThumbnailURL: thumbnailURL,
	}

	return media, nil
}

func (p *FapodropProvider) GetRecentMediaID(username string) (int, error) {
	c := colly.NewCollector()

	userPageURL, err := url.JoinPath(p.BaseURL, username)
	if err != nil {
		return 0, err
	}

	isFound := false
	recentMediaID := 0

	c.OnHTML(fmt.Sprintf(".one-pack a[href^='/%s']", username), func(e *colly.HTMLElement) {
		if !isFound {
			href := e.Attr("href")

			mediaID, err := p.parseMediaID(href)
			if err != nil {
				return
			}

			recentMediaID = mediaID
			isFound = true
		}
	})

	err = c.Visit(userPageURL)
	if err != nil {
		return 0, fmt.Errorf("failed to visit %s: %w", userPageURL, err)
	}

	return recentMediaID, nil
}

func (p *FapodropProvider) buildURL(baseURL string, name string) (string, error) {
	firstSymbol := name[0]
	secondSymbol := name[1]

	mediaURL, err := url.JoinPath(baseURL, "images", string(firstSymbol), string(secondSymbol), name, "1", "photo")
	if err != nil {
		return "", err
	}

	return mediaURL, nil
}

func (p *FapodropProvider) buildThumbnailURL(baseURL string, name string) (string, error) {
	firstSymbol := name[0]
	secondSymbol := name[1]

	mediaThumbnailURL, err := url.JoinPath(baseURL, "images", string(firstSymbol), string(secondSymbol), name, "1", "thumbnails")
	if err != nil {
		return "", err
	}

	return mediaThumbnailURL, nil
}

func (p *FapodropProvider) parseMediaID(url string) (int, error) {
	re := regexp.MustCompile(`\/\d{4}$`)

	match := re.FindString(url)
	if match == "" {
		return 0, fmt.Errorf("invalid url: %s", url)
	}

	numStr := match[1:]                    // Take out the first "/"
	numStr = strings.TrimLeft(numStr, "0") // Remove leading zeros

	mediaID, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}

	return mediaID, err
}
