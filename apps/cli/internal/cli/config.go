package cli

import (
	"fmt"

	"github.com/gentij/taskforge/apps/cli/internal/config"
)

func loadConfig() (config.Config, string, error) {
	path := config.ResolvePath(configPath)
	cfg, err := config.Load(path)
	if err != nil {
		return config.Config{}, path, err
	}

	if cfg.ServerURL == "" {
		cfg.ServerURL = serverURL
	}

	return cfg, path, nil
}

func saveConfig(path string, cfg config.Config) error {
	if cfg.ServerURL == "" {
		return fmt.Errorf("server URL is required")
	}

	return config.Save(path, cfg)
}
