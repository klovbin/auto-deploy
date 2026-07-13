package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func repoNameFromURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimSuffix(raw, ".git")
	raw = strings.TrimSuffix(raw, "/")

	if raw == "" {
		return "", fmt.Errorf("ссылка на репозиторий пустая")
	}

	if strings.HasPrefix(raw, "git@") {
		if idx := strings.LastIndex(raw, ":"); idx >= 0 {
			raw = raw[idx+1:]
		}
	} else {
		raw = strings.TrimPrefix(raw, "https://")
		raw = strings.TrimPrefix(raw, "http://")
		raw = strings.TrimPrefix(raw, "ssh://")

		if idx := strings.Index(raw, "/"); idx >= 0 {
			raw = raw[idx+1:]
		}
	}

	parts := strings.Split(raw, "/")
	name := parts[len(parts)-1]
	if name == "" {
		return "", fmt.Errorf("не удалось определить имя репозитория")
	}

	return name, nil
}

func generateDeployKeys(cfg Config) error {
	if cfg.Repository == "" {
		return fmt.Errorf("сначала добавьте репозиторий")
	}

	repoName, err := repoNameFromURL(cfg.Repository)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	deployDir := filepath.Join(cwd, repoName, ".deploy")
	privateKeyPath := filepath.Join(deployDir, "id_ed25519")
	publicKeyPath := filepath.Join(deployDir, "id_ed25519.pub")

	if err := generateSSHKeyPair(privateKeyPath, publicKeyPath); err != nil {
		return err
	}

	fmt.Printf("Ключи созданы в %s\n", deployDir)
	fmt.Printf("  приватный: %s\n", privateKeyPath)
	fmt.Printf("  публичный: %s\n", publicKeyPath)

	return nil
}
