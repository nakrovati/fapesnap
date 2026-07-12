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
	cfg        *config.Config
	cancel     context.CancelFunc
	mu         sync.Mutex
}

func NewAppService(app *application.App, cfg *config.Config) *AppService {
	return &AppService{
		app: app,
		cfg: cfg,
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
		return []providers.Media{}, fmt.Errorf("failed to resolve collection slug: %w", err)
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
		return fmt.Errorf("failed to resolve collection slug: %w", err)
	}

	mediaItems, err := a.scraper.GetMediaItems(collectionSlug)
	if err != nil {
		return err
	}

	for i := range mediaItems {
		if mediaItems[i].URL == "" {
			mediaURL, err := a.scraper.ResolveMediaURL(mediaItems[i].PageURL, collectionSlug)
			if err != nil {
				return fmt.Errorf("failed to resolve media URL: %w", err)
			}

			mediaItems[i].URL = mediaURL
		}
	}

	err = a.downloader.DownloadMediaItems(ctx, mediaItems, a.cfg.DownloadDir, providerName, collectionSlug, maxParallelDownloads)
	if err != nil {
		a.app.Logger.Error("Error downloading media files", "error", err)

		return err
	}

	a.app.Logger.Info("All media downloaded successfully")

	return nil
}

func (a *AppService) DownloadMedia(pageURL string, collectionInput string, providerName string) error {
	a.StopTask()

	scr, err := scraper.NewScraper(providerName)
	if err != nil {
		return err
	}

	a.scraper = scr

	collectionSlug, err := a.scraper.ResolveCollectionSlug(collectionInput)
	if err != nil {
		return fmt.Errorf("failed to resolve collection slug: %w", err)
	}

	downloadDir, err := utils.GetCollectionDownloadDir(a.cfg.DownloadDir, providerName, collectionSlug)
	if err != nil {
		return fmt.Errorf("failed to get download directory: %w", err)
	}

	mediaURL, err := a.scraper.ResolveMediaURL(pageURL, collectionSlug)
	if err != nil {
		return err
	}

	err = a.downloader.DownloadMedia(a.app.Context(), pageURL, mediaURL, downloadDir)
	if err != nil {
		return fmt.Errorf("error downloading media: %w", err)
	}

	return nil
}

func (a *AppService) GetDownloadDir() config.DownloadDir {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.cfg.DownloadDir
}

func (a *AppService) SelectDownloadDir() (*config.DownloadDir, error) {
	path, err := a.app.Dialog.OpenFile().SetTitle("Select a download directory").CanChooseDirectories(true).CanChooseFiles(false).PromptForSingleSelection()
	if err != nil {
		return &config.DownloadDir{}, err
	}

	err = utils.ValidateDir(path)
	if err != nil {
		return &config.DownloadDir{}, err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.cfg.DownloadDir = config.DownloadDir{
		Path:      path,
		IsDefault: false,
	}

	err = a.cfg.Save()
	if err != nil {
		return nil, err
	}

	return &a.cfg.DownloadDir, nil
}

func (a *AppService) UnsetDownloadDir() (*config.DownloadDir, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	cfg, err := config.Default()
	if err != nil {
		return nil, err
	}

	a.cfg.DownloadDir = cfg.DownloadDir

	err = a.cfg.Save()
	if err != nil {
		return nil, err
	}

	return &a.cfg.DownloadDir, nil
}

func (a *AppService) StopTask() {
	if a.cancel != nil {
		a.cancel()
		a.cancel = nil
	}
}

//nolint:unparam
func (a *AppService) ServiceStartup(_ctx context.Context, _options application.ServiceOptions) error {
	a.downloader = downloader.NewDownloader(a.app.Logger, a.app.Event)

	return nil
}
