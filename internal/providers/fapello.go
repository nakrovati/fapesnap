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

func (p *FapelloProvider) GetCollectionFromURL(inputURL string) (string, error) {
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

func (p *FapelloProvider) GetMedia(mediaID string, username string) (Media, error) {
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
		URL:          mediaURL,
		ThumbnailURL: thumbnailURL,
	}

	return media, nil
}

func (p *FapelloProvider) GetRecentMediaID(username string) (int, error) {
	c := colly.NewCollector()

	userSrc, err := url.JoinPath(p.BaseURL, username)
	if err != nil {
		return 0, err
	}

	isFound := false
	recentMediaID := 0

	c.OnHTML(fmt.Sprintf("#content div a[href*='%s']", username), func(e *colly.HTMLElement) {
		if !isFound {
			src := e.Attr("href")

			mediaID, err := p.parseMediaID(src)
			if err != nil {
				return
			}

			recentMediaID = mediaID
			isFound = true
		}
	})

	err = c.Visit(userSrc)
	if err != nil {
		return 0, fmt.Errorf("failed to visit %s: %w", userSrc, err)
	}

	if !isFound {
		return 0, errors.New("user not found")
	}

	return recentMediaID, nil
}

func (p *FapelloProvider) buildURL(baseURL string, username string, recentID int) (string, error) {
	firstSymbol := string(username[0])
	secondSymbol := string(username[1])
	countGroup := strconv.Itoa(utils.RoundUp(recentID))

	mediaURL, err := url.JoinPath(baseURL, "content", firstSymbol, secondSymbol, username, countGroup)
	if err != nil {
		return "", err
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
