package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nakrovati/fapesnap/internal/config"
)

func GetCollectionDownloadDir(
	downloadDir config.DownloadDir,
	providerName string,
	collectionSlug string,
) (string, error) {
	baseDir := downloadDir.Path

	if downloadDir.IsDefault {
		baseDir = filepath.Join(baseDir, "fapesnap")
	}

	collectionDir := filepath.Join(baseDir, providerName, collectionSlug)

	if err := os.MkdirAll(collectionDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create download directory: %w", err)
	}

	return collectionDir, nil
}
