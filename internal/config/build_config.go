package config

import (
	"embed"
	"fmt"

	"go.yaml.in/yaml/v3"
)

func LoadBuildConfig(buildConfig embed.FS) (*BuildConfig, error) {
	data, err := buildConfig.ReadFile("build/config.yml")
	if err != nil {
		return &BuildConfig{}, fmt.Errorf("read config: %w", err)
	}

	var cfg *BuildConfig

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return &BuildConfig{}, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}
