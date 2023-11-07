package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	computev1 "github.com/mohaijiang/computeshare-client/api/compute/v1"
	"github.com/mohaijiang/computeshare-client/internal/conf"
	"github.com/mohaijiang/computeshare-client/internal/service"
	"github.com/mohaijiang/computeshare-client/third_party/agent"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	vmDockerService *service.VmDockerService,
	computePowerService *service.ComputePowerService,
	agentService *agent.AgentService,
	vmWebsocketHandler *service.VmWebsocketHandler,
	cronJob *service.CronJob,
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
	computev1.RegisterVmDockerHTTPServer(srv, vmDockerService)
	computev1.RegisterComputePowerClientHTTPServer(srv, computePowerService)

	srv.HandleFunc("/v1/vm/{id}/terminal", vmWebsocketHandler.Terminal)

	// register
	err := agentService.Register()
	if err != nil {
		panic(err)
	}

	cronJob.StartJob()

	return srv
}
