package main

import (
	"context"
	"fmt"

	"github.com/nakrovati/fapesnap/internal/downloader"
	"github.com/nakrovati/fapesnap/internal/scraper"
)

// App struct.
type App struct {
	scraper    *scraper.Scraper
	downloader downloader.Downloader
	ctx        context.Context
}

// NewApp creates a new App application struct.
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetPhotos(collection string, provider string) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("Collection cannot be empty")
	}

	a.scraper = scraper.NewScraper(provider)
	a.scraper.SetContext(a.ctx)

	photos, err := a.scraper.GetPhotoURLs(collection)
	if err != nil {
		return nil, err
	}

	return photos, nil
}

func (a *App) DownloadPhotos(collection string, provider string) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("Collection cannot be empty")
	}

	a.scraper = scraper.NewScraper(provider)

	photos, err := a.scraper.GetPhotoURLs(collection)
	if err != nil {
		return nil, err
	}

	a.downloader = downloader.Downloader{}
	a.downloader.SetContext(a.ctx)

	err = a.downloader.DownloadPhotos(photos, provider, collection)
	if err != nil {
		return nil, fmt.Errorf("Error downloading photos: %v", err)
	}
	return photos, nil
}
