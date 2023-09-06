package server

import (
	computev1 "computeshare-client/api/compute/v1"
	v1 "computeshare-client/api/helloworld/v1"
	p2pv1 "computeshare-client/api/network/v1"
	"computeshare-client/internal/conf"
	"computeshare-client/internal/service"
	"github.com/go-kratos/swagger-api/openapiv2"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	greeter *service.GreeterService,
	p2pService *service.P2pService,
	vmService *service.VmService,
	computepowerService *service.ComputepowerService,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	openAPIhandler := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", openAPIhandler)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	p2pv1.RegisterP2PHTTPServer(srv, p2pService)
	computev1.RegisterVmHTTPServer(srv, vmService)
	computev1.RegisterComputepowerHTTPServer(srv, computepowerService)
	return srv
}