package configstore

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/trust/deploy/internal/domain/config"
)

const fileName = ".trust-deploy.json"

func path() (string, error) {
	return filepath.Join(config.DefaultWorkDir, fileName), nil
}

func Load() (config.Config, error) {
	configPath, err := path()
	if err != nil {
		return config.Config{}, err
	}

	cfg, err := readConfig(configPath)
	if err == nil {
		return cfg, nil
	}
	if !os.IsNotExist(err) {
		return config.Config{}, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return config.Config{}, nil
	}

	legacyPath := filepath.Join(cwd, fileName)
	cfg, legacyErr := readConfig(legacyPath)
	if legacyErr != nil {
		if os.IsNotExist(legacyErr) {
			return config.Config{}, nil
		}
		return config.Config{}, legacyErr
	}

	if saveErr := Save(cfg); saveErr == nil {
		_ = os.Remove(legacyPath)
	}

	return cfg, nil
}

func readConfig(configPath string) (config.Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
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

	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0o644)
}
