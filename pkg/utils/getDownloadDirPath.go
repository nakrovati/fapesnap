package utils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func GetDownloadDirectory(providerName string, userName string) (string, error) {
	usr, err := user.Current()

	if err != nil {
		return "", err
	}
	downloadDir := filepath.Join(usr.HomeDir, "Downloads", "fapesnap", providerName, userName)

	err = os.MkdirAll(downloadDir, 0755)
	if err != nil {
		fmt.Println("Error while creating download directory:", err)
		return "", err
	}

	return downloadDir, nil
}
