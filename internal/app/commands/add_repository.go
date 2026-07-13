package commands

import (
	"strings"

	"github.com/trust/deploy/internal/domain/config"
	configstore "github.com/trust/deploy/internal/infra/out/config"
)

type AddRepositoryHandler struct{}

func (AddRepositoryHandler) Handle(cfg *config.Config, url string) error {
	cfg.Repository = strings.TrimSpace(url)
	return configstore.Save(*cfg)
}
