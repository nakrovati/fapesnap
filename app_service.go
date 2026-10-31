package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/egoist/mygo"
	"github.com/nakrovati/fapesnap/internal/config"
	"github.com/nakrovati/fapesnap/internal/downloader"
	"github.com/nakrovati/fapesnap/internal/pkg/utils"
	"github.com/nakrovati/fapesnap/internal/providers"
	"github.com/nakrovati/fapesnap/internal/scraper"
)

type AppService struct {
	downloader  *downloader.Downloader
	cfg         *config.Config
	taskMu      sync.Mutex
	currentTask *appTask
	configMu    *sync.RWMutex
}

type appTask struct {
	cancel context.CancelFunc
}

func NewAppService(cfg *config.Config, dl *downloader.Downloader, configMu *sync.RWMutex) *AppService {
	return &AppService{
		cfg:        cfg,
		downloader: dl,
		configMu:   configMu,
	}
}

func (a *AppService) GetMediaItems(_ context.Context, collectionInput string, providerName string) ([]providers.Media, error) {
	a.StopTask()

	scr, err := scraper.NewScraper(providerName)
	if err != nil {
		return []providers.Media{}, err
	}

	collectionSlug, err := scr.ResolveCollectionSlug(collectionInput)
	if err != nil {
		return []providers.Media{}, fmt.Errorf("failed to resolve collection slug: %w", err)
	}

	mediaItems, err := scr.GetMediaItems(collectionSlug)
	if err != nil {
		return nil, err
	}

	return mediaItems, nil
}

func (a *AppService) DownloadMediaItems(ctx context.Context, collectionInput string, providerName string, maxParallelDownloads int) error {
	ctx, task := a.startTask(ctx)
	defer a.finishTask(task)

	scr, err := scraper.NewScraper(providerName)
	if err != nil {
		return err
	}

	collectionSlug, err := scr.ResolveCollectionSlug(collectionInput)
	if err != nil {
		return fmt.Errorf("failed to resolve collection slug: %w", err)
	}

	mediaItems, err := scr.GetMediaItems(collectionSlug)
	if err != nil {
		return err
	}

	for i := range mediaItems {
		if mediaItems[i].URL == "" {
			mediaURL, err := scr.ResolveMediaURL(mediaItems[i].PageURL, collectionSlug)
			if err != nil {
				return fmt.Errorf("failed to resolve media URL: %w", err)
			}

			mediaItems[i].URL = mediaURL
		}
	}

	a.configMu.RLock()
	downloadDir := a.cfg.DownloadDir
	a.configMu.RUnlock()

	err = a.downloader.DownloadMediaItems(ctx, mediaItems, downloadDir, providerName, collectionSlug, maxParallelDownloads)
	if err != nil {
		slog.Error("Error downloading media files", "error", err)

		return err
	}

	slog.Info("All media downloaded successfully")

	return nil
}

func (a *AppService) DownloadMedia(ctx context.Context, pageURL string, collectionInput string, providerName string) error {
	ctx, task := a.startTask(ctx)
	defer a.finishTask(task)

	scr, err := scraper.NewScraper(providerName)
	if err != nil {
		return err
	}

	collectionSlug, err := scr.ResolveCollectionSlug(collectionInput)
	if err != nil {
		return fmt.Errorf("failed to resolve collection slug: %w", err)
	}

	a.configMu.RLock()
	baseDownloadDir := a.cfg.DownloadDir
	a.configMu.RUnlock()

	downloadDir, err := utils.GetCollectionDownloadDir(baseDownloadDir, providerName, collectionSlug)
	if err != nil {
		return fmt.Errorf("failed to get download directory: %w", err)
	}

	mediaURL, err := scr.ResolveMediaURL(pageURL, collectionSlug)
	if err != nil {
		return err
	}

	err = a.downloader.DownloadMedia(ctx, pageURL, mediaURL, downloadDir)
	if err != nil {
		return fmt.Errorf("error downloading media: %w", err)
	}

	return nil
}

func (a *AppService) GetDownloadDir() config.DownloadDir {
	a.configMu.RLock()
	defer a.configMu.RUnlock()

	return a.cfg.DownloadDir
}

func (a *AppService) SelectDownloadDir(ctx context.Context) (*config.DownloadDir, error) {
	paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
		Parent:    mygo.CallerWindow(ctx),
		Title:     "Select a download directory",
		Directory: true,
	})
	if err != nil {
		return nil, err
	}

	if len(paths) == 0 {
		dir := a.GetDownloadDir()

		return &dir, nil
	}

	err = utils.ValidateDir(paths[0])
	if err != nil {
		return nil, err
	}

	a.configMu.Lock()
	defer a.configMu.Unlock()

	a.cfg.DownloadDir = config.DownloadDir{
		Path:      paths[0],
		IsDefault: false,
	}

	err = a.cfg.Save()
	if err != nil {
		return nil, err
	}

	dir := a.cfg.DownloadDir

	return &dir, nil
}

func (a *AppService) UnsetDownloadDir() (*config.DownloadDir, error) {
	a.configMu.Lock()
	defer a.configMu.Unlock()

	cfg, err := config.Default()
	if err != nil {
		return nil, err
	}

	a.cfg.DownloadDir = cfg.DownloadDir

	err = a.cfg.Save()
	if err != nil {
		return nil, err
	}

	dir := a.cfg.DownloadDir

	return &dir, nil
}

func (a *AppService) StopTask() {
	a.taskMu.Lock()
	defer a.taskMu.Unlock()

	if a.currentTask != nil {
		a.currentTask.cancel()
		a.currentTask = nil
	}
}

func (a *AppService) OpenExternal(url string) error {
	return mygo.Shell.OpenExternal(url)
}

func (a *AppService) startTask(parent context.Context) (context.Context, *appTask) {
	a.taskMu.Lock()
	defer a.taskMu.Unlock()

	if a.currentTask != nil {
		a.currentTask.cancel()
	}

	ctx, cancel := context.WithCancel(parent)
	task := &appTask{cancel: cancel}
	a.currentTask = task

	return ctx, task
}

func (a *AppService) finishTask(task *appTask) {
	a.taskMu.Lock()
	defer a.taskMu.Unlock()

	task.cancel()

	if a.currentTask == task {
		a.currentTask = nil
	}
}
