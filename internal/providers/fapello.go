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
	MaxPhotoID   int
	MinPhotoID   int
	ProviderName string
	BaseURL      string
}

func NewFapelloProvider() *FapelloProvider {
	return &FapelloProvider{
		MaxPhotoID:   100000,
		MinPhotoID:   1,
		ProviderName: "fapello",
		BaseURL:      "https://fapello.com",
	}
}

func (p *FapelloProvider) FetchPhotoURLs(collectionSlug string) ([]Photo, error) {
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

func (p *FapelloProvider) GetPhoto(photoID string, username string) (Photo, error) {
	intPhotoID, err := strconv.Atoi(photoID)
	if err != nil {
		return Photo{}, err
	}

	paddedID := fmt.Sprintf("%04d", intPhotoID)
	photoName := fmt.Sprintf("%s_%v.jpg", username, paddedID)
	photoThumbnailName := fmt.Sprintf("%s_%v_300px.jpg", username, paddedID)

	urlWithoutID, err := p.buildURL(p.BaseURL, username, intPhotoID)
	if err != nil {
		return Photo{}, err
	}

	photoURL, err := url.JoinPath(urlWithoutID, photoName)
	if err != nil {
		return Photo{}, err
	}

	thumbnailURL, err := url.JoinPath(urlWithoutID, photoThumbnailName)
	if err != nil {
		return Photo{}, err
	}

	photo := Photo{
		URL:          photoURL,
		ThumbnailURL: thumbnailURL,
	}

	return photo, nil
}

func (p *FapelloProvider) GetRecentPhotoID(username string) (int, error) {
	c := colly.NewCollector()

	userSrc, err := url.JoinPath(p.BaseURL, username)
	if err != nil {
		return 0, err
	}

	isFound := false
	recentPhotoID := 0

	c.OnHTML(fmt.Sprintf("#content div a[href*='%s']", username), func(e *colly.HTMLElement) {
		if !isFound {
			src := e.Attr("href")

			photoID, err := p.parsePhotoID(src)
			if err != nil {
				return
			}

			recentPhotoID = photoID
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

	return recentPhotoID, nil
}

func (p *FapelloProvider) buildURL(baseURL string, username string, recentID int) (string, error) {
	firstSymbol := string(username[0])
	secondSymbol := string(username[1])
	photoCountGroup := strconv.Itoa(utils.RoundUp(recentID))

	photoURL, err := url.JoinPath(baseURL, "content", firstSymbol, secondSymbol, username, photoCountGroup)
	if err != nil {
		return "", err
	}

	return photoURL, nil
}

func (p *FapelloProvider) parsePhotoID(url string) (int, error) {
	re := regexp.MustCompile(`\/(\d+)/$`)

	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return 0, fmt.Errorf("invalid url: %s", url)
	}

	photoID, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	return photoID, nil
}
