//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"computeshare-client/internal/biz"
	"computeshare-client/internal/conf"
	"computeshare-client/internal/data"
	"computeshare-client/internal/server"
	"computeshare-client/internal/service"
	"computeshare-client/third_party/agent"
	"computeshare-client/third_party/p2p"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(
		wire.Build(
			server.ProviderSet,
			data.ProviderSet,
			biz.ProviderSet,
			service.ProviderSet,
			p2p.ProviderSet,
			agent.ProviderSet,
			newApp,
		),
	)
}
