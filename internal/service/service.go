package service

import (
	"github.com/docker/docker/client"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewGreeterService,
	NewP2pService,
	NewDockerCli,
	NewVmService,
	NewComputePowerService,
	NewVmWebsocketHandler,
)

func NewDockerCli() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv)
}
