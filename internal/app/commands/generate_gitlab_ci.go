package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/trust/deploy/internal/domain/config"
	"github.com/trust/deploy/internal/shared/paths"
)

const gitlabCIFileName = ".gitlab-ci.yml"

func GitLabCIExists(cfg config.Config) bool {
	repoPath, err := paths.RepoDirectory(cfg)
	if err != nil {
		return false
	}

	_, err = os.Stat(filepath.Join(repoPath, gitlabCIFileName))
	return err == nil
}

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

type GenerateGitLabCIHandler struct{}

func (GenerateGitLabCIHandler) Handle(cfg config.Config) error {
	repoPath, repoName, err := paths.RequireClonedRepo(cfg)
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
