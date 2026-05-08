package utils

import (
	"fmt"
	"os"
)

func ValidateDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("not a directory")
	}

	return nil
}
