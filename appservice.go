package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/nakrovati/fapesnap/internal/config"
	"github.com/nakrovati/fapesnap/internal/downloader"
	"github.com/nakrovati/fapesnap/internal/pkg/utils"
	"github.com/nakrovati/fapesnap/internal/providers"
	"github.com/nakrovati/fapesnap/internal/scraper"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type AppService struct {
	app        *application.App
	scraper    *scraper.Scraper
	downloader *downloader.Downloader
	config     *config.Config
	cancel     context.CancelFunc
	mu         sync.RWMutex
}

func NewAppService(app *application.App) *AppService {
	return &AppService{
		app: app,
	}
}

func (a *AppService) GetPhotos(collectionInput string, providerName string) ([]providers.Photo, error) {
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

func (a *AppService) DownloadPhotos(collectionInput string, providerName string, maxParallelDownloads int) error {
	a.StopTask()

	if collectionInput == "" {
		return errors.New("collection cannot be empty")
	}

	ctx, cancel := context.WithCancel(a.app.Context())
	a.cancel = cancel

	a.scraper = scraper.NewScraper(providerName)

	collectionSlug, err := a.scraper.ResolveCollectionSlug(collectionInput)
	if err != nil {
		return fmt.Errorf("failed to resolve collection slug: %w", err)
	}

	photoURLs, err := a.scraper.GetPhotoURLs(collectionSlug)
	if err != nil {
		return err
	}

	err = a.downloader.DownloadPhotos(a.app, ctx, photoURLs, a.config.DownloadDir, providerName, collectionSlug, maxParallelDownloads)
	if err != nil {
		fmt.Printf("Error downloading photos: %v\n", err)
	} else {
		fmt.Println("All photos downloaded successfully.")
	}

	return nil
}

func (a *AppService) DownloadPhoto(src string, collectionInput string, providerName string) error {
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

	downloadDir, err := utils.GetCollectionDownloadDir(a.config.DownloadDir, providerName, collectionSlug)
	if err != nil {
		return fmt.Errorf("failed to get download directory: %w", err)
	}

	err = a.downloader.DownloadPhoto(a.app.Context(), src, downloadDir)
	if err != nil {
		return fmt.Errorf("error downloading photo: %w", err)
	}

	return nil
}

func (a *AppService) GetDownloadDir() config.DownloadDir {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.config.DownloadDir
}

func (a *AppService) SelectDownloadDir() (*config.DownloadDir, error) {
	path, err := a.app.Dialog.OpenFile().SetTitle("Select a download directory").CanChooseDirectories(true).CanChooseFiles(false).PromptForSingleSelection()
	if err != nil {
		return &config.DownloadDir{}, err
	}

	if err := utils.ValidateDir(path); err != nil {
		return &config.DownloadDir{}, err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.config.DownloadDir = config.DownloadDir{
		Path:      path,
		IsDefault: false,
	}

	if err := config.Save(a.config); err != nil {
		return nil, err
	}

	return &a.config.DownloadDir, nil
}

func (a *AppService) UnsetDownloadDir() (*config.DownloadDir, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.config.DownloadDir = config.Default().DownloadDir

	if err := config.Save(a.config); err != nil {
		return nil, err
	}

	return &a.config.DownloadDir, nil
}

func (a *AppService) StopTask() {
	if a.cancel != nil {
		a.cancel()
		a.cancel = nil
	}
}

func (a *AppService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	a.downloader = downloader.NewDownloader()

	cfg, err := config.Load()
	if err != nil {
		a.app.Logger.Error(
			"Error occurred", "failed to load config",
			err,
		)

		a.config = config.Default()

		return nil
	}

	a.config = cfg

	return nil
}
