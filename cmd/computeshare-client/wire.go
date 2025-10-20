//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/tntlinking-computeshare/computeshare-client/internal/biz"
	"github.com/tntlinking-computeshare/computeshare-client/internal/conf"
	"github.com/tntlinking-computeshare/computeshare-client/internal/server"
	"github.com/tntlinking-computeshare/computeshare-client/internal/service"
	"github.com/tntlinking-computeshare/computeshare-client/third_party/agent"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(
		wire.Build(
			server.ProviderSet,
			//data.ProviderSet,
			biz.ProviderSet,
			service.ProviderSet,
			agent.ProviderSet,
			newApp,
		),
	)
}
