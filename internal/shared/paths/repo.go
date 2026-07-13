package paths

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/trust/deploy/internal/domain/config"
	"github.com/trust/deploy/internal/domain/repository"
)

func RepoDirectory(cfg config.Config) (string, error) {
	if cfg.Repository == "" {
		return "", fmt.Errorf("сначала добавьте репозиторий")
	}

	repoName, err := repository.NameFromURL(cfg.Repository)
	if err != nil {
		return "", err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(cwd, repoName), nil
}

func RequireClonedRepo(cfg config.Config) (repoPath string, repoName string, err error) {
	repoPath, err = RepoDirectory(cfg)
	if err != nil {
		return "", "", err
	}

	if _, err := os.Stat(filepath.Join(repoPath, ".git")); err != nil {
		return "", "", fmt.Errorf("сначала склонируйте репозиторий")
	}

	repoName, err = repository.NameFromURL(cfg.Repository)
	if err != nil {
		return "", "", err
	}

	return repoPath, repoName, nil
}
