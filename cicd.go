package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const gitlabCIFileName = ".gitlab-ci.yml"

func gitlabCITemplate(repoName string) string {
	return fmt.Sprintf(`stages:
  - deploy

deploy:
  stage: deploy
  image: alpine:3.21
  before_script:
    - apk add --no-cache openssh-client
    - eval "$(ssh-agent -s)"
    - echo "$SSH_PRIVATE_KEY" | tr -d '\r' | ssh-add -
    - mkdir -p ~/.ssh && chmod 700 ~/.ssh
    - ssh-keyscan -p 22 "$SSH_HOST" >> ~/.ssh/known_hosts
  script:
    - |
      ssh -p 22 "root@${SSH_HOST}" "
        set -e
        cd /var/www/%s
        git pull
        docker compose up -d --build
      "
  rules:
    - if: $CI_COMMIT_BRANCH == "main"
`, repoName)
}

func requireClonedRepo(cfg Config) (string, string, error) {
	repoPath, err := repoDirectory(cfg)
	if err != nil {
		return "", "", err
	}

	if !isGitRepo(repoPath) {
		return "", "", fmt.Errorf("сначала склонируйте репозиторий")
	}

	repoName, err := repoNameFromURL(cfg.Repository)
	if err != nil {
		return "", "", err
	}

	return repoPath, repoName, nil
}

func generateGitLabCI(cfg Config) error {
	repoPath, repoName, err := requireClonedRepo(cfg)
	if err != nil {
		return err
	}

	ciPath := filepath.Join(repoPath, gitlabCIFileName)
	if _, err := os.Stat(ciPath); err == nil {
		fmt.Printf("%s уже существует\n", gitlabCIFileName)
		return nil
	}

	content := gitlabCITemplate(repoName)
	if err := os.WriteFile(ciPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("не удалось создать %s: %w", gitlabCIFileName, err)
	}

	fmt.Printf("%s создан: %s\n", gitlabCIFileName, ciPath)
	return nil
}

func generateCICDKey(cfg Config) error {
	repoPath, _, err := requireClonedRepo(cfg)
	if err != nil {
		return err
	}

	deployDir := filepath.Join(repoPath, ".deploy")
	privateKeyPath := filepath.Join(deployDir, "SSH_PRIVATE_KEY")
	publicKeyPath := filepath.Join(deployDir, "SSH_PRIVATE_KEY.pub")

	if err := generateSSHKeyPair(privateKeyPath, publicKeyPath); err != nil {
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
