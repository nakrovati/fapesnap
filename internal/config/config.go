package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"time"
)

//nolint:gochecknoglobals
var mu sync.RWMutex

type DownloadDir struct {
	Path      string `json:"path"`
	IsDefault bool   `json:"isDefault"`
}

type CheckInterval int

const (
	IntervalDay   CheckInterval = 1440
	Interval3Days CheckInterval = 4320
	IntervalWeek  CheckInterval = 10080
	IntervalMonth CheckInterval = 43200
)

func (i CheckInterval) IsValid() bool {
	switch i {
	case IntervalDay, Interval3Days, IntervalWeek, IntervalMonth:
		return true
	}

	return false
}

type Updates struct {
	AutoCheck            bool          `json:"autoCheck"`
	CheckOnStartup       bool          `json:"checkOnStartup"`
	CheckIntervalMinutes CheckInterval `json:"checkIntervalMinutes"`
	IncludePrereleases   bool          `json:"includePrereleases"`
	LastCheckedAt        *time.Time    `json:"lastCheckedAt"`
}

type Config struct {
	DownloadDir DownloadDir `json:"downloadDir"`
	Updates     Updates     `json:"updates"`
}

type Info struct {
	Version string `yaml:"version"`
}

type BuildConfig struct {
	Info Info `yaml:"info"`
}

const (
	defaultAutoCheck            = true
	defaultCheckOnStartup       = false
	defaultCheckIntervalMinutes = IntervalWeek
	defaultPrereleaseEnabled    = false
)

func Default() (*Config, error) {
	downloadDir, err := defaultDownloadDir()
	if err != nil {
		return nil, err
	}

	return &Config{
		DownloadDir: downloadDir,
		Updates:     defaultUpdates(),
	}, nil
}

func Load() (*Config, error) {
	mu.Lock()
	defer mu.Unlock()

	path, err := configPath()
	if err != nil {
		return nil, err
	}

	// #nosec G304
	data, err := os.ReadFile(path)

	if os.IsNotExist(err) {
		cfg, err := Default()
		if err != nil {
			return nil, err
		}

		err = saveUnlocked(cfg)
		if err != nil {
			return nil, err
		}

		return cfg, nil
	}

	if err != nil {
		return nil, err
	}

	return parseConfig(data)
}

func (c *Config) Save() error {
	mu.Lock()
	defer mu.Unlock()

	return saveUnlocked(c)
}

func configPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "fapesnap", "config.json"), nil
}

func saveUnlocked(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), 0o750)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := path + ".tmp"

	err = os.WriteFile(tmpPath, data, 0o600)
	if err != nil {
		return err
	}

	return os.Rename(tmpPath, path)
}

func parseConfig(data []byte) (*Config, error) {
	cfg, err := Default()
	if err != nil {
		return nil, err
	}

	var raw struct {
		DownloadDir *DownloadDir `json:"downloadDir"`
		Updates     *Updates     `json:"updates"`
	}

	err = json.Unmarshal(data, &raw)
	if err != nil {
		return nil, err
	}

	if raw.DownloadDir != nil {
		cfg.DownloadDir = *raw.DownloadDir
	}

	if raw.Updates != nil {
		cfg.Updates = *raw.Updates
	}

	normalize(cfg)

	return cfg, nil
}

func normalize(cfg *Config) {
	if cfg.DownloadDir.Path == "" {
		cfg.DownloadDir, _ = defaultDownloadDir()
	}

	if cfg.Updates.CheckIntervalMinutes <= 0 {
		cfg.Updates.CheckIntervalMinutes = defaultCheckIntervalMinutes
	}
}

func defaultDownloadDir() (DownloadDir, error) {
	usr, err := user.Current()
	if err != nil {
		return DownloadDir{}, err
	}

	return DownloadDir{
		Path:      filepath.Join(usr.HomeDir, "Downloads"),
		IsDefault: true,
	}, nil
}

func defaultUpdates() Updates {
	return Updates{
		AutoCheck:            defaultAutoCheck,
		CheckOnStartup:       defaultCheckOnStartup,
		CheckIntervalMinutes: defaultCheckIntervalMinutes,
		IncludePrereleases:   defaultPrereleaseEnabled,
	}
}
