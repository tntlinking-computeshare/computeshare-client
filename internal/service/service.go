package service

import (
	"github.com/docker/docker/client"
	"github.com/google/wire"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/mohaijiang/computeshare-client/internal/conf"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewGreeterService,
	NewDockerCli,
	NewVmService,
	NewComputePowerService,
	NewVmWebsocketHandler,
	NewIpfShell,
)

func NewDockerCli() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv)
}

func NewIpfShell(c *conf.Data) *shell.Shell {
	return shell.NewShell(c.Ipfs.Url)
}
