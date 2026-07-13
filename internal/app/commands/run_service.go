package commands

import (
	"github.com/trust/deploy/internal/domain/config"
	"github.com/trust/deploy/internal/infra/out/docker"
	"github.com/trust/deploy/internal/shared/paths"
)

type RunServiceHandler struct{}

func (RunServiceHandler) Handle(cfg config.Config) error {
	repoPath, _, err := paths.RequireClonedRepo(cfg)
	if err != nil {
		return err
	}

	return docker.RunService(repoPath)
}

func CanRunService(cfg config.Config) bool {
	repoPath, _, err := paths.RequireClonedRepo(cfg)
	if err != nil {
		return false
	}
	return docker.ComposeFileExists(repoPath)
}
