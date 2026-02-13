package config

import (
	"path/filepath"

	"os"
)

func DefaultConfigPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "config.json"
	}

	return filepath.Join(dir, "taskforge", "config.json")
}
