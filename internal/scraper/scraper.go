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

func (s *Scraper) GetPhotoURLs(collection string) ([]string, error) {
	_, err := url.Parse(collection)
	if err != nil {
		collectionString, err := s.provider.GetCollectionFromURL(collection)
		if err != nil {
			return []string{}, err
		}

		collection = collectionString
	}

	return s.provider.FetchPhotoURLs(collection)
}
