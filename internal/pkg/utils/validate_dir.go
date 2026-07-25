package utils

import (
	"fmt"
	"os"
)

func ValidateDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat directory %q: %w", path, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path %q is not a directory", path)
	}

	return nil
}
