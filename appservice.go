package main

import (
	"context"
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

func (a *AppService) GetMediaItems(collectionInput string, providerName string) ([]providers.Media, error) {
	a.StopTask()

	scr, err := scraper.NewScraper(providerName)
	if err != nil {
		return []providers.Media{}, err
	}

	a.scraper = scr

	collectionSlug, err := a.scraper.ResolveCollectionSlug(collectionInput)
	if err != nil {
		return []providers.Media{}, fmt.Errorf("Failed to resolve collection slug: %w", err)
	}

	mediaItems, err := a.scraper.GetMediaItems(collectionSlug)
	if err != nil {
		return nil, err
	}

	return mediaItems, nil
}

func (a *AppService) DownloadMediaItems(collectionInput string, providerName string, maxParallelDownloads int) error {
	a.StopTask()

	ctx, cancel := context.WithCancel(a.app.Context())
	a.cancel = cancel

	scr, err := scraper.NewScraper(providerName)
	if err != nil {
		return err
	}

	a.scraper = scr

	collectionSlug, err := a.scraper.ResolveCollectionSlug(collectionInput)
	if err != nil {
		return fmt.Errorf("Failed to resolve collection slug: %w", err)
	}

	mediaItems, err := a.scraper.GetMediaItems(collectionSlug)
	if err != nil {
		return err
	}

	err = a.downloader.DownloadMediaItems(ctx, mediaItems, a.config.DownloadDir, providerName, collectionSlug, maxParallelDownloads)
	if err != nil {
		a.app.Logger.Error("Error downloading media files", "error", err)
	} else {
		a.app.Logger.Info("All media downloaded successfully")
	}

	return nil
}

func (a *AppService) DownloadMedia(src string, collectionInput string, providerName string) error {
	a.StopTask()

	scr, err := scraper.NewScraper(providerName)
	if err != nil {
		return err
	}

	a.scraper = scr

	collectionSlug, err := a.scraper.ResolveCollectionSlug(collectionInput)
	if err != nil {
		return fmt.Errorf("Failed to resolve collection slug: %w", err)
	}

	downloadDir, err := utils.GetCollectionDownloadDir(a.config.DownloadDir, providerName, collectionSlug)
	if err != nil {
		return fmt.Errorf("Failed to get download directory: %w", err)
	}

	err = a.downloader.DownloadMedia(a.app.Context(), src, downloadDir)
	if err != nil {
		return fmt.Errorf("Error downloading media: %w", err)
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
	a.downloader = downloader.NewDownloader(a.app.Logger, a.app.Event)

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
