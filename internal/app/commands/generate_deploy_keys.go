package commands

import (
	"fmt"
	"path/filepath"

	"github.com/trust/deploy/internal/domain/config"
	"github.com/trust/deploy/internal/pkg/sshkeys"
	"github.com/trust/deploy/internal/shared/paths"
)

type GenerateDeployKeysHandler struct{}

func (GenerateDeployKeysHandler) Handle(cfg config.Config) error {
	if cfg.Repository == "" {
		return fmt.Errorf("сначала добавьте репозиторий")
	}

	deployDir, err := paths.DeployDirectory(cfg)
	if err != nil {
		return err
	}

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
