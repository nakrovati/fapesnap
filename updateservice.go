package main

import (
	"context"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/nakrovati/fapesnap/internal/config"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/updater"
	"github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

type UpdateService struct {
	app      *application.App
	cfg      *config.Config
	buildCfg *config.BuildConfig
	mu       sync.RWMutex

	//nolint:containedctx
	periodicCtx    context.Context
	periodicCancel context.CancelFunc
	periodicDone   chan struct{}
}

func NewUpdateService(app *application.App, cfg *config.Config, buildCfg *config.BuildConfig) *UpdateService {
	return &UpdateService{
		app:      app,
		cfg:      cfg,
		buildCfg: buildCfg,
	}
}

func (u *UpdateService) Check() {
	u.checkNow()

	u.mu.Lock()
	defer u.mu.Unlock()

	if u.cfg.Updates.AutoCheck {
		u.stopPeriodicCheckLocked()
		u.startPeriodicCheckLocked()
	}
}

func (u *UpdateService) GetUpdatesConfig() config.Updates {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.cfg.Updates
}

func (u *UpdateService) SetUpdatesConfig(updates config.Updates) (config.Updates, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	prereleaseChanged := u.cfg.Updates.IncludePrereleases != updates.IncludePrereleases
	intervalChanged := u.cfg.Updates.CheckIntervalMinutes != updates.CheckIntervalMinutes
	autoCheckChanged := u.cfg.Updates.AutoCheck != updates.AutoCheck

	if !updates.CheckIntervalMinutes.IsValid() {
		updates.CheckIntervalMinutes = config.IntervalWeek
	}

	u.cfg.Updates = updates

	err := u.cfg.Save()
	if err != nil {
		return config.Updates{}, err
	}

	if prereleaseChanged {
		_ = restartApp(u.app)

		return u.cfg.Updates, nil
	}

	if intervalChanged || autoCheckChanged {
		u.stopPeriodicCheckLocked()

		if u.cfg.Updates.AutoCheck {
			u.startPeriodicCheckLocked()
		}
	}

	return u.cfg.Updates, nil
}

func (u *UpdateService) ServiceStartup(_ctx context.Context, _options application.ServiceOptions) error {
	gh, err := github.New(github.Config{
		Repository: "nakrovati/fapesnap",
		Prerelease: u.cfg.Updates.IncludePrereleases,
	})
	if err != nil {
		return err
	}

	err = u.app.Updater.Init(updater.Config{
		CurrentVersion: u.buildCfg.Info.Version,
		Providers:      []updater.Provider{gh},
	})
	if err != nil {
		return err
	}

	if u.cfg.Updates.CheckOnStartup {
		//nolint:contextcheck
		u.checkNow()
	}

	u.mu.Lock()
	if u.cfg.Updates.AutoCheck {
		//nolint:contextcheck
		u.startPeriodicCheckLocked()
	}
	u.mu.Unlock()

	return nil
}

func (u *UpdateService) startPeriodicCheckLocked() {
	if u.periodicCancel != nil {
		return
	}

	interval := time.Duration(u.cfg.Updates.CheckIntervalMinutes) * time.Minute
	if interval <= 0 {
		interval = 24 * time.Hour
	}

	u.periodicCtx, u.periodicCancel = context.WithCancel(context.Background())
	u.periodicDone = make(chan struct{})

	go u.periodicCheckLoop(interval)
}

func (u *UpdateService) stopPeriodicCheckLocked() {
	if u.periodicCancel != nil {
		u.periodicCancel()
		<-u.periodicDone
		u.periodicCancel = nil
		u.periodicDone = nil
	}
}

func (u *UpdateService) periodicCheckLoop(interval time.Duration) {
	defer close(u.periodicDone)

	initialDelay := interval

	u.mu.RLock()
	lastChecked := u.cfg.Updates.LastCheckedAt
	u.mu.RUnlock()

	if lastChecked != nil {
		elapsed := time.Since(*lastChecked)
		if elapsed < interval {
			initialDelay = interval - elapsed
		} else {
			initialDelay = 0
		}
	}

	t := time.NewTimer(initialDelay)
	defer t.Stop()

	for {
		select {
		case <-u.periodicCtx.Done():
			return
		case <-t.C:
			t.Reset(interval)

			s := u.app.Updater.State()
			if s == updater.StateChecking || s == updater.StateDownloading || s == updater.StateVerifying || s == updater.StateInstalling {
				continue
			}

			u.checkNow()
		}
	}
}

func (u *UpdateService) checkNow() {
	err := u.app.Updater.CheckAndInstall(context.Background())
	if err != nil {
		u.app.Logger.Error("Update check failed", "error", err)
	} else {
		u.app.Logger.Info("Update check completed", "state", u.app.Updater.State())
	}

	u.updateLastCheckedAt()
}

func (u *UpdateService) updateLastCheckedAt() {
	u.mu.Lock()
	defer u.mu.Unlock()

	now := time.Now().UTC()
	u.cfg.Updates.LastCheckedAt = &now

	err := u.cfg.Save()
	if err != nil {
		u.app.Logger.Error("Failed to save update config", "error", err)
	}
}

func restartApp(app *application.App) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.Command(exe, os.Args[1:]...)

	err = cmd.Start()
	if err != nil {
		return err
	}

	app.Quit()

	return nil
}
