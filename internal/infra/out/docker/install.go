package docker

import (
	"fmt"
	"os"
	"os/exec"
)

func IsInstalled() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

func Install() error {
	if IsInstalled() {
		fmt.Println("Docker уже установлен")
		return nil
	}

	fmt.Println("Устанавливаю Docker...")
	cmd := exec.Command("sh", "-c", "curl -fsSL https://get.docker.com | sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("установка Docker: %w", err)
	}

	fmt.Println("Docker установлен")
	return nil
}
