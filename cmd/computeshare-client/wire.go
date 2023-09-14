//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/mohaijiang/computeshare-client/internal/biz"
	"github.com/mohaijiang/computeshare-client/internal/conf"
	"github.com/mohaijiang/computeshare-client/internal/data"
	"github.com/mohaijiang/computeshare-client/internal/server"
	"github.com/mohaijiang/computeshare-client/internal/service"
	"github.com/mohaijiang/computeshare-client/third_party/agent"
	"github.com/mohaijiang/computeshare-client/third_party/p2p"

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
