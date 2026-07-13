package commands

import (
	"github.com/trust/deploy/internal/domain/config"
	"github.com/trust/deploy/internal/infra/out/git"
)

type CloneRepositoryHandler struct{}

func (CloneRepositoryHandler) Handle(cfg config.Config) error {
	return git.Clone(cfg)
}
