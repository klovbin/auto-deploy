package commands

import "github.com/trust/deploy/internal/infra/out/docker"

type InstallDockerHandler struct{}

func (InstallDockerHandler) Handle() error {
	return docker.Install()
}
