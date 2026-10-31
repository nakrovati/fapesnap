package main

import (
	"context"
	"sync"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/plugins/updater"
	"github.com/nakrovati/fapesnap/internal/config"
)

type UpdateService struct {
	cfg      *config.Config
	configMu *sync.RWMutex
}

func NewUpdateService(cfg *config.Config, configMu *sync.RWMutex) *UpdateService {
	return &UpdateService{
		cfg:      cfg,
		configMu: configMu,
	}
}

func (u *UpdateService) GetUpdatesConfig() config.Updates {
	u.configMu.RLock()
	updates := u.cfg.Updates
	u.configMu.RUnlock()

	updates.AutoCheck = updater.AutomaticChecks()

	return updates
}

func (u *UpdateService) UpdatesEnabled() bool {
	return mygo.Updater.Enabled()
}

func (u *UpdateService) SetUpdatesConfig(updates config.Updates) (config.Updates, error) {
	u.configMu.Lock()
	defer u.configMu.Unlock()

	if !updates.CheckIntervalMinutes.IsValid() {
		updates.CheckIntervalMinutes = config.IntervalWeek
	}

	previous := updater.AutomaticChecks()

	u.cfg.Updates = updates
	if err := u.cfg.Save(); err != nil {
		return config.Updates{}, err
	}

	mygo.Updater.SetPrereleases(updates.IncludePrereleases)

	if previous != updates.AutoCheck {
		updater.SetAutomaticChecks(updates.AutoCheck)
	}

	updates.AutoCheck = updater.AutomaticChecks()

	return updates, nil
}

func (u *UpdateService) Check(context.Context) {
	updater.CheckForUpdates()
}
