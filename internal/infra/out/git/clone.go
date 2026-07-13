package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/trust/deploy/internal/domain/config"
	"github.com/trust/deploy/internal/shared/paths"
)

func isGitRepo(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil
}

func env(repoPath string) []string {
	gitEnv := os.Environ()
	privateKey := filepath.Join(repoPath, ".deploy", "id_ed25519")
	if _, err := os.Stat(privateKey); err == nil {
		sshCmd := fmt.Sprintf("ssh -i %q -o IdentitiesOnly=yes -o StrictHostKeyChecking=accept-new", privateKey)
		gitEnv = append(gitEnv, "GIT_SSH_COMMAND="+sshCmd)
	}
	return gitEnv
}

func run(dir string, gitEnv []string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = gitEnv
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func detectDefaultBranch(repoPath string) (string, error) {
	out, err := exec.Command("git", "-C", repoPath, "branch", "-r").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("не удалось получить список веток: %w", err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "origin/HEAD -> ") {
			return strings.TrimPrefix(line, "origin/HEAD -> origin/"), nil
		}
	}

	for _, branch := range []string{"main", "master"} {
		for _, line := range strings.Split(string(out), "\n") {
			if strings.TrimSpace(line) == "origin/"+branch {
				return branch, nil
			}
		}
	}

	return "", fmt.Errorf("не удалось определить основную ветку")
}

func ensureGitignore(repoPath string) error {
	const entry = ".deploy"

	gitignorePath := filepath.Join(repoPath, ".gitignore")
	content, err := os.ReadFile(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("не удалось прочитать .gitignore: %w", err)
	}

	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == entry || trimmed == entry+"/" {
			fmt.Println(".deploy уже есть в .gitignore")
			return nil
		}
	}

	var builder strings.Builder
	if len(content) > 0 {
		builder.Write(content)
		if !strings.HasSuffix(string(content), "\n") {
			builder.WriteByte('\n')
		}
	}
	builder.WriteString(entry)
	builder.WriteByte('\n')

	if err := os.WriteFile(gitignorePath, []byte(builder.String()), 0o644); err != nil {
		return fmt.Errorf("не удалось обновить .gitignore: %w", err)
	}

	fmt.Println(".deploy добавлен в .gitignore")
	return nil
}

func Clone(cfg config.Config) error {
	repoPath, err := paths.RepoDirectory(cfg)
	if err != nil {
		return err
	}

	if isGitRepo(repoPath) {
		return fmt.Errorf("репозиторий уже склонирован: %s", repoPath)
	}

	if err := os.MkdirAll(repoPath, 0o755); err != nil {
		return fmt.Errorf("не удалось создать папку: %w", err)
	}

	gitEnv := env(repoPath)

	if err := run(repoPath, gitEnv, "init"); err != nil {
		return fmt.Errorf("git init: %w", err)
	}

	if err := run(repoPath, gitEnv, "remote", "add", "origin", cfg.Repository); err != nil {
		return fmt.Errorf("git remote add: %w", err)
	}

	fmt.Println("Клонирование...")
	if err := run(repoPath, gitEnv, "fetch", "origin"); err != nil {
		return fmt.Errorf("git fetch: %w", err)
	}

	branch, err := detectDefaultBranch(repoPath)
	if err != nil {
		return err
	}

	if err := run(repoPath, gitEnv, "checkout", "-t", "origin/"+branch); err != nil {
		return fmt.Errorf("git checkout: %w", err)
	}

	if err := ensureGitignore(repoPath); err != nil {
		return err
	}

	fmt.Printf("Репозиторий склонирован в %s\n", repoPath)
	return nil
}
