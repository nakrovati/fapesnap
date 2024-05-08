package utils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func GetDownloadDirectory(providerName string, collection string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	downloadDir := filepath.Join(usr.HomeDir, "Downloads", "fapesnap", providerName, collection)

	err = os.MkdirAll(downloadDir, 0o755)
	if err != nil {
		return "", fmt.Errorf("failed to create download directory: %w", err)
	}

	return downloadDir, nil
}
