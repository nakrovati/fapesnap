package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
)

type DownloadDir struct {
	Path      string `json:"path"`
	IsDefault bool   `json:"isDefault"`
}

type Config struct {
	DownloadDir DownloadDir `json:"downloadDir"`
}

func Default() *Config {
	return &Config{
		DownloadDir: defaultDownloadDir(),
	}
}

func Load() (*Config, error) {
	path := configPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := Default()

			if err := Save(cfg); err != nil {
				return nil, err
			}

			return cfg, nil
		}

		return nil, err
	}

	cfg := Default()

	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.DownloadDir.Path == "" {
		cfg.DownloadDir = defaultDownloadDir()
	}

	return cfg, nil
}

func Save(cfg *Config) error {
	path := configPath()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

func configPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(configDir, "fapesnap", "config.json")
}

func defaultDownloadDir() DownloadDir {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return DownloadDir{
		Path:      filepath.Join(usr.HomeDir, "Downloads"),
		IsDefault: true,
	}
}
