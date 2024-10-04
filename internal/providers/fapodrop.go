package providers

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"slices"
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

func (p *FapodropProvider) InitProvider() {
	p.ProviderName = "fapodrop"
	p.BaseURL = "https://fapodrop.com"
}

func (p FapodropProvider) FetchPhotoURLs(collection string) ([]string, error) {
	if p.MinPhotoID > p.MaxPhotoID {
		return []string{}, fmt.Errorf("min photo ID (%d) is greater than max photo ID (%d)", p.MinPhotoID, p.MaxPhotoID)
	}

	recentPhotoID, err := p.GetRecentPhotoID(collection)
	if err != nil {
		return []string{}, err
	}

	if p.MaxPhotoID > recentPhotoID {
		p.MaxPhotoID = recentPhotoID
	}

	photos := make([]string, 0, p.MaxPhotoID-p.MinPhotoID+1)

	for i := p.MinPhotoID; i <= p.MaxPhotoID; i++ {
		photoID := strconv.Itoa(i)

		photoURL, err := p.GetPhotoURL(photoID, collection)
		if err != nil {
			return []string{}, err
		}

		photos = append(photos, photoURL)
	}

	slices.Reverse(photos)

	return photos, nil
}

func (p FapodropProvider) GetCollectionFromURL(inputURL string) (string, error) {
	_, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	if !strings.Contains(inputURL, p.BaseURL) {
		return "", errors.New("Unvalid domain")
	}

	inputURL = strings.TrimSuffix(inputURL, "/")
	parts := strings.Split(inputURL, "/")

	if len(parts) < 4 || parts[len(parts)-1] == "" {
		return "", errors.New("Can't get collection from url")
	}

	return parts[len(parts)-1], nil
}

func (p FapodropProvider) GetPhotoURL(photoID string, username string) (string, error) {
	intPhotoID, err := strconv.Atoi(photoID)
	if err != nil {
		return "", err
	}

	urlWithoutID, err := p.buildURL(p.BaseURL, username)
	if err != nil {
		return "", err
	}

	paddedID := fmt.Sprintf("%04d", intPhotoID)
	photoName := fmt.Sprintf("%s_%s.jpeg", username, paddedID)

	url, err := url.JoinPath(urlWithoutID, photoName)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (p FapodropProvider) GetRecentPhotoID(username string) (int, error) {
	c := colly.NewCollector()

	userSrc, err := url.JoinPath(p.BaseURL, username)
	if err != nil {
		return 0, err
	}

	isFound := false
	recentPhotoID := 0

	c.OnHTML(fmt.Sprintf(".row .one-pack a[href^='/%s']", username), func(e *colly.HTMLElement) {
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

	return recentPhotoID, nil
}

func (p FapodropProvider) buildURL(baseURL string, name string) (string, error) {
	firstSymbol := name[0]
	secondSymbol := name[1]

	photoURL, err := url.JoinPath(baseURL, "images", string(firstSymbol), string(secondSymbol), name, "1", "photo")
	if err != nil {
		return "", err
	}

	return photoURL, nil
}

func (p FapodropProvider) parsePhotoID(url string) (int, error) {
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
