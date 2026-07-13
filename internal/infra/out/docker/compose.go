package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ComposeFileExists(repoPath string) bool {
	if _, err := os.Stat(filepath.Join(repoPath, "docker-compose.yml")); err == nil {
		return true
	}
	if _, err := os.Stat(filepath.Join(repoPath, "docker-compose.yaml")); err == nil {
		return true
	}
	return false
}

func RunService(repoPath string) error {
	if !IsInstalled() {
		return fmt.Errorf("сначала установите Docker")
	}

	if !ComposeFileExists(repoPath) {
		return fmt.Errorf("docker-compose.yml не найден в %s", repoPath)
	}

	fmt.Printf("Запускаю сервис в %s...\n", repoPath)
	cmd := exec.Command("docker", "compose", "up", "-d", "--build")
	cmd.Dir = repoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker compose up: %w", err)
	}

	fmt.Println("Сервис запущен")
	return nil
}
