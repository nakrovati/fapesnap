package providers

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/nakrovati/fapesnap/internal/pkg/utils"

	"github.com/gocolly/colly/v2"
)

type FapelloProvider struct {
	MaxPhotoID   int
	MinPhotoID   int
	ProviderName string
	BaseURL      string
}

func (p *FapelloProvider) InitProvider() {
	p.ProviderName = "fapello"
	p.BaseURL = "https://fapello.com"
}

func (p FapelloProvider) FetchPhotoURLs(collection string) ([]string, error) {
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

	photos := make([]string, p.MaxPhotoID-p.MinPhotoID+1)

	for i := p.MinPhotoID; i <= p.MaxPhotoID; i++ {
		photoID := strconv.Itoa(i)
		photos[i-p.MinPhotoID], err = p.GetPhotoURL(photoID, collection)
		if err != nil {
			return []string{}, err
		}
	}

	return photos, nil
}

func (p FapelloProvider) GetCollectionFromURL(inputURL string) (string, error) {
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

func (p FapelloProvider) GetPhotoURL(photoID string, username string) (string, error) {
	intPhotoID, err := strconv.Atoi(photoID)
	if err != nil {
		return "", err
	}

	urlWithoutID, err := p.buildURL(p.BaseURL, username, intPhotoID)
	if err != nil {
		return "", err
	}

	paddedID := fmt.Sprintf("%04d", intPhotoID)
	photoName := fmt.Sprintf("%s_%v.jpg", username, paddedID)

	url, err := url.JoinPath(urlWithoutID, photoName)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (p FapelloProvider) GetRecentPhotoID(username string) (int, error) {
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
		return 0, fmt.Errorf("user not found")
	}

	return recentPhotoID, nil
}

func (p FapelloProvider) buildURL(baseURL string, username string, recentID int) (string, error) {
	firstSymbol := string(username[0])
	secondSymbol := string(username[1])
	photoCountGroup := strconv.Itoa(utils.RoundUp(recentID))

	photoURL, err := url.JoinPath(baseURL, "content", firstSymbol, secondSymbol, username, photoCountGroup)
	if err != nil {
		return "", err
	}

	return photoURL, nil
}

func (p FapelloProvider) parsePhotoID(url string) (int, error) {
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
