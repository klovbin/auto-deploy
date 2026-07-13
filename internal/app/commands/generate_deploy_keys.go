package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/trust/deploy/internal/domain/config"
	"github.com/trust/deploy/internal/domain/repository"
	"github.com/trust/deploy/internal/pkg/sshkeys"
)

type GenerateDeployKeysHandler struct{}

func (GenerateDeployKeysHandler) Handle(cfg config.Config) error {
	if cfg.Repository == "" {
		return fmt.Errorf("сначала добавьте репозиторий")
	}

	repoName, err := repository.NameFromURL(cfg.Repository)
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

	if err := sshkeys.GeneratePair(privateKeyPath, publicKeyPath); err != nil {
		return err
	}

	fmt.Printf("Ключи созданы в %s\n", deployDir)
	fmt.Printf("  приватный: %s\n", privateKeyPath)
	fmt.Printf("  публичный: %s\n", publicKeyPath)

	return nil
}
