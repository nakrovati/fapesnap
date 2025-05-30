package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/nakrovati/fapesnap/internal/downloader"
	"github.com/nakrovati/fapesnap/internal/scraper"
)

// App struct.
type App struct {
	scraper    *scraper.Scraper
	downloader downloader.Downloader
	//nolint:containedctx
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewApp creates a new App application struct.
func NewApp() *App {
	return &App{}
}

func (a *App) GetPhotos(collection string, provider string) ([]string, error) {
	a.StopTask()

	if collection == "" {
		return nil, errors.New("collection cannot be empty")
	}

	a.scraper = scraper.NewScraper(provider)

	photos, err := a.scraper.GetPhotoURLs(collection)
	if err != nil {
		return nil, fmt.Errorf("failed to get photo URLs: %w", err)
	}

	return photos, nil
}

func (a *App) DownloadPhotos(collectionName string, providerName string, maxParallelDownloads int) error {
	a.StopTask()

	if collectionName == "" {
		return errors.New("collection cannot be empty")
	}

	ctx, cancel := context.WithCancel(a.ctx)
	a.cancelFunc = cancel

	a.scraper = scraper.NewScraper(providerName)

	photoURLs, err := a.scraper.GetPhotoURLs(collectionName)
	if err != nil {
		return err
	}

	err = a.downloader.DownloadPhotos(ctx, photoURLs, providerName, collectionName, maxParallelDownloads)
	if err != nil {
		fmt.Printf("Error downloading photos: %v\n", err)
	} else {
		fmt.Println("All photos downloaded successfully.")
	}

	return nil
}

func (a *App) StopTask() {
	if a.cancelFunc != nil {
		a.cancelFunc()
		a.cancelFunc = nil
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.downloader = downloader.Downloader{}
}
