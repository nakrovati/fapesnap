package main

import (
	"log"
	"sync"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/plugins/updater"
	"github.com/nakrovati/fapesnap/internal/config"
	"github.com/nakrovati/fapesnap/internal/downloader"
)

var configMu sync.RWMutex

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	appService := NewAppService(cfg, downloader.NewDownloader(nil), &configMu)
	updateService := NewUpdateService(cfg, &configMu)
	mygo.Use(updater.New(updater.Options{
		Interval:               time.Duration(cfg.Updates.CheckIntervalMinutes) * time.Minute,
		DisableAutomaticChecks: !cfg.Updates.AutoCheck,
	}))

	mygo.Updater.SetPrereleases(cfg.Updates.IncludePrereleases)
	mygo.Bind(appService, updateService)

	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:           "Fapesnap",
			Width:           960,
			Height:          640,
			MinWidth:        400,
			MinHeight:       300,
			BackgroundColor: "light-dark(#ffffff, #0a0a0a)", // the page's --background in src/routes/layout.css
			StateKey:        "main",
			URL:             "/",
		})
	})

	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
