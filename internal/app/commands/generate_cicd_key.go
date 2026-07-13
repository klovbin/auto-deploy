package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/trust/deploy/internal/domain/config"
	"github.com/trust/deploy/internal/pkg/sshkeys"
	"github.com/trust/deploy/internal/shared/paths"
)

type GenerateCICDKeyHandler struct{}

func (GenerateCICDKeyHandler) Handle(cfg config.Config) error {
	repoPath, _, err := paths.RequireClonedRepo(cfg)
	if err != nil {
		return err
	}

	deployDir := filepath.Join(repoPath, ".deploy")
	privateKeyPath := filepath.Join(deployDir, "SSH_PRIVATE_KEY")
	publicKeyPath := filepath.Join(deployDir, "SSH_PRIVATE_KEY.pub")

	if err := sshkeys.GeneratePair(privateKeyPath, publicKeyPath); err != nil {
		return err
	}

	if err := os.Chmod(privateKeyPath, 0o600); err != nil {
		return fmt.Errorf("не удалось выставить права на приватный ключ: %w", err)
	}

	fmt.Printf("CI/CD ключ создан в %s\n", deployDir)
	fmt.Printf("  приватный: %s (600)\n", privateKeyPath)
	fmt.Printf("  публичный: %s\n", publicKeyPath)
	fmt.Println()
	fmt.Println("Добавьте содержимое SSH_PRIVATE_KEY в переменную GitLab CI: SSH_PRIVATE_KEY")

	return nil
}
