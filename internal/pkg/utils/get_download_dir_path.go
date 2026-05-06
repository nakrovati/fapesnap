package utils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func GetDownloadDirectory(providerName string, collectionSlug string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	downloadDir := filepath.Join(usr.HomeDir, "Downloads", "fapesnap", providerName, collectionSlug)

	err = os.MkdirAll(downloadDir, 0750)
	if err != nil {
		return "", fmt.Errorf("failed to create download directory: %w", err)
	}

	return downloadDir, nil
}
