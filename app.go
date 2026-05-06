package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/nakrovati/fapesnap/internal/downloader"
	"github.com/nakrovati/fapesnap/internal/pkg/utils"
	"github.com/nakrovati/fapesnap/internal/providers"
	"github.com/nakrovati/fapesnap/internal/scraper"
)

// App struct.
type App struct {
	scraper    *scraper.Scraper
	downloader *downloader.Downloader
	//nolint:containedctx
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewApp creates a new App application struct.
func NewApp() *App {
	return &App{}
}

func (a *App) GetPhotos(collectionInput string, providerName string) ([]providers.Photo, error) {
	a.StopTask()

	if collectionInput == "" {
		return nil, errors.New("collection cannot be empty")
	}

	a.scraper = scraper.NewScraper(providerName)

	collectionSlug, err := a.scraper.ResolveCollectionSlug(collectionInput)
	if err != nil {
		return []providers.Photo{}, fmt.Errorf("failed to resolve collection slug: %w", err)
	}

	photos, err := a.scraper.GetPhotoURLs(collectionSlug)
	if err != nil {
		return nil, fmt.Errorf("failed to get photo URLs: %w", err)
	}

	return photos, nil
}

func (a *App) DownloadPhotos(collectionInput string, providerName string, maxParallelDownloads int) error {
	a.StopTask()

	if collectionInput == "" {
		return errors.New("collection cannot be empty")
	}

	ctx, cancel := context.WithCancel(a.ctx)
	a.cancelFunc = cancel

	a.scraper = scraper.NewScraper(providerName)

	collectionSlug, err := a.scraper.ResolveCollectionSlug(collectionInput)
	if err != nil {
		return fmt.Errorf("failed to resolve collection slug: %w", err)
	}

	photoURLs, err := a.scraper.GetPhotoURLs(collectionSlug)
	if err != nil {
		return err
	}

	err = a.downloader.DownloadPhotos(ctx, photoURLs, providerName, collectionSlug, maxParallelDownloads)
	if err != nil {
		fmt.Printf("Error downloading photos: %v\n", err)
	} else {
		fmt.Println("All photos downloaded successfully.")
	}

	return nil
}

func (a *App) DownloadPhoto(src string, collectionInput string, providerName string) error {
	a.StopTask()

	provider := providers.GetProvider(providerName)
	if provider == nil {
		return nil
	}

	s := scraper.NewScraper(providerName)

	collectionSlug, err := s.ResolveCollectionSlug(collectionInput)
	if err != nil {
		return fmt.Errorf("failed to resolve collection slug: %w", err)
	}

	downloadDir, err := utils.GetDownloadDirectory(providerName, collectionSlug)
	if err != nil {
		return fmt.Errorf("failed to get download directory: %w", err)
	}

	err = a.downloader.DownloadPhoto(a.ctx, src, downloadDir)
	if err != nil {
		return fmt.Errorf("error downloading photo: %w", err)
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
	a.downloader = downloader.NewDownloader()
}
