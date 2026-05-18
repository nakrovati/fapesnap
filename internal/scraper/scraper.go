package scraper

import (
	"errors"
	"net/url"
	"strings"

	"github.com/nakrovati/fapesnap/internal/providers"
)

type Scraper struct {
	provider providers.Provider
}

var ErrInvalidProvider = errors.New("invalid provider")

func NewScraper(providerName string) (*Scraper, error) {
	provider := providers.GetProvider(providerName)
	if provider == nil {
		return nil, ErrInvalidProvider
	}

	return &Scraper{
		provider: provider,
	}, nil
}

func (s *Scraper) GetMediaItems(collectionSlug string) ([]providers.Media, error) {
	return s.provider.FetchMediaItems(collectionSlug)
}

func (s *Scraper) ResolveCollectionSlug(collectionInput string) (string, error) {
	collectionInput = strings.TrimSpace(collectionInput)
	if collectionInput == "" {
		return "", errors.New("collection cannot be empty")
	}

	u, err := url.Parse(collectionInput)
	if err != nil {
		return "", err
	}

	if u.Scheme != "" && u.Host != "" {
		return s.provider.GetCollectionFromURL(collectionInput)
	}

	return collectionInput, nil
}
