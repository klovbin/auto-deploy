package configstore

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/trust/deploy/internal/domain/config"
)

const fileName = ".trust-deploy.json"

func path() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, fileName), nil
}

func Load() (config.Config, error) {
	configPath, err := path()
	if err != nil {
		return config.Config{}, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return config.Config{}, nil
		}
		return config.Config{}, err
	}

	var cfg config.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

func Save(cfg config.Config) error {
	configPath, err := path()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0o644)
}
