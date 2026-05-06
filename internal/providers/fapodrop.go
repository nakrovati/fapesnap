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
	MaxPhotoID   int
	MinPhotoID   int
	ProviderName string
	BaseURL      string
}

func NewFapodropProvider() *FapodropProvider {
	return &FapodropProvider{
		MaxPhotoID:   100000,
		MinPhotoID:   1,
		ProviderName: "fapodrop",
		BaseURL:      "https://fapodrop.com",
	}
}

func (p *FapodropProvider) FetchPhotoURLs(collectionSlug string) ([]Photo, error) {
	if p.MinPhotoID > p.MaxPhotoID {
		return []Photo{}, fmt.Errorf("min photo ID (%d) is greater than max photo ID (%d)", p.MinPhotoID, p.MaxPhotoID)
	}

	recentPhotoID, err := p.GetRecentPhotoID(collectionSlug)
	if err != nil {
		return []Photo{}, err
	}

	if p.MaxPhotoID > recentPhotoID {
		p.MaxPhotoID = recentPhotoID
	}

	photos := make([]Photo, 0, p.MaxPhotoID-p.MinPhotoID+1)

	for i := p.MaxPhotoID; i >= p.MinPhotoID; i-- {
		photo, err := p.GetPhoto(strconv.Itoa(i), collectionSlug)
		if err != nil {
			fmt.Printf("Failed to get photo: %v\n", err)

			continue
		}

		photos = append(photos, photo)
	}

	if len(photos) == 0 {
		return []Photo{}, errors.New("no photos found")
	}

	return photos, nil
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

func (p *FapodropProvider) GetPhoto(photoID string, username string) (Photo, error) {
	intPhotoID, err := strconv.Atoi(photoID)
	if err != nil {
		return Photo{}, err
	}

	paddedID := fmt.Sprintf("%04d", intPhotoID)
	photoName := fmt.Sprintf("%s_%s.jpeg", username, paddedID)
	photoThumbnailName := fmt.Sprintf("%s_%s_thumbnail.jpeg", username, paddedID)

	urlWithoutID, err := p.buildURL(p.BaseURL, username)
	if err != nil {
		return Photo{}, err
	}

	thumbnailURLWithoutID, err := p.buildThumbnailURL(p.BaseURL, username)
	if err != nil {
		return Photo{}, err
	}

	photoURL, err := url.JoinPath(urlWithoutID, photoName)
	if err != nil {
		return Photo{}, err
	}

	thumbnailURL, err := url.JoinPath(thumbnailURLWithoutID, photoThumbnailName)
	if err != nil {
		return Photo{}, err
	}

	photo := Photo{
		URL:          photoURL,
		ThumbnailURL: thumbnailURL,
	}

	return photo, nil
}

func (p *FapodropProvider) GetRecentPhotoID(username string) (int, error) {
	c := colly.NewCollector()

	userPageURL, err := url.JoinPath(p.BaseURL, username)
	if err != nil {
		return 0, err
	}

	isFound := false
	recentPhotoID := 0

	c.OnHTML(fmt.Sprintf(".one-pack a[href^='/%s']", username), func(e *colly.HTMLElement) {
		if !isFound {
			href := e.Attr("href")

			photoID, err := p.parsePhotoID(href)
			if err != nil {
				return
			}

			recentPhotoID = photoID
			isFound = true
		}
	})

	err = c.Visit(userPageURL)
	if err != nil {
		return 0, fmt.Errorf("failed to visit %s: %w", userPageURL, err)
	}

	return recentPhotoID, nil
}

func (p *FapodropProvider) buildURL(baseURL string, name string) (string, error) {
	firstSymbol := name[0]
	secondSymbol := name[1]

	photoURL, err := url.JoinPath(baseURL, "images", string(firstSymbol), string(secondSymbol), name, "1", "photo")
	if err != nil {
		return "", err
	}

	return photoURL, nil
}

func (p *FapodropProvider) buildThumbnailURL(baseURL string, name string) (string, error) {
	firstSymbol := name[0]
	secondSymbol := name[1]

	photoURL, err := url.JoinPath(baseURL, "images", string(firstSymbol), string(secondSymbol), name, "1", "thumbnails")
	if err != nil {
		return "", err
	}

	return photoURL, nil
}

func (p *FapodropProvider) parsePhotoID(url string) (int, error) {
	re := regexp.MustCompile(`\/\d{4}$`)

	match := re.FindString(url)
	if match == "" {
		return 0, fmt.Errorf("invalid url: %s", url)
	}

	numStr := match[1:]                    // Take out the first "/"
	numStr = strings.TrimLeft(numStr, "0") // Убираем ведущие нули

	photoID, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}

	return photoID, err
}
