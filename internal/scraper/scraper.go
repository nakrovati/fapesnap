package scraper

import (
	"net/url"

	"github.com/nakrovati/fapesnap/internal/providers"
)

type Scraper struct {
	provider providers.Provider
}

func NewScraper(providerName string) *Scraper {
	provider := providers.GetProvider(providerName)
	if provider == nil {
		return nil
	}

	return &Scraper{
		provider: provider,
	}
}

func (s *Scraper) GetPhotoURLs(collectionSlug string) ([]providers.Photo, error) {
	return s.provider.FetchPhotoURLs(collectionSlug)
}

func (s *Scraper) ResolveCollectionSlug(collectionInput string) (string, error) {
	u, err := url.Parse(collectionInput)
	if err == nil && u.Scheme != "" && u.Host != "" {
		collectionSlug, err := s.provider.GetCollectionFromURL(collectionInput)
		if err != nil {
			return "", err
		}

		return collectionSlug, nil
	}

	return collectionInput, nil
}
