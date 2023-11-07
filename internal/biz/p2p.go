package biz

import (
	"context"
	"github.com/fatedier/frp/client"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/mohaijiang/computeshare-client/internal/conf"
	"github.com/samber/lo"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type P2pClient struct {
	svr         *client.Service
	cfg         *v1.ClientCommonConfig
	pxyCfgs     []v1.ProxyConfigurer
	visitorCfgs []v1.VisitorConfigurer
}

func NewP2pClient(c *conf.Server) (*P2pClient, error) {

	cfg, pxyCfgs, visitorCfgs, _, err := generate(c.P2P.GatewayIp, int(c.P2P.GatewayPort))
	if err != nil {
		return nil, err
	}

	svr, err := client.NewService(cfg, pxyCfgs, visitorCfgs, "")
	shouldGracefulClose := cfg.Transport.Protocol == "kcp" || cfg.Transport.Protocol == "quic"
	if shouldGracefulClose {
		go handleTermSignal(svr)
	}

	go func() {
		err = svr.Run(context.Background())
		if err != nil {
			panic(err)
		}
	}()
	return &P2pClient{
		svr:         svr,
		cfg:         cfg,
		pxyCfgs:     pxyCfgs,
		visitorCfgs: visitorCfgs,
	}, err
}

func (c *P2pClient) CreateProxy(name string, localIp string, localPort, remotePort int64) (string, int, error) {

	proxyConfigurer := &v1.TCPProxyConfig{
		ProxyBaseConfig: v1.ProxyBaseConfig{
			Name: name,
			Type: "tcp",
			ProxyBackend: v1.ProxyBackend{
				LocalIP:   localIp,
				LocalPort: int(localPort),
			},
		},
		RemotePort: int(remotePort),
	}
	proxyConfigurer.Complete("")

	c.pxyCfgs = append(c.pxyCfgs, proxyConfigurer)

	err := c.svr.ReloadConf(c.pxyCfgs, c.visitorCfgs)

	return c.cfg.ServerAddr, int(remotePort), err
}

func (c *P2pClient) DeleteProxy(name string) error {
	c.pxyCfgs = lo.FilterMap(c.pxyCfgs, func(item v1.ProxyConfigurer, _ int) (v1.ProxyConfigurer, bool) {
		return item, item.GetBaseConfig().Name != name
	})

	return c.svr.ReloadConf(c.pxyCfgs, c.visitorCfgs)
}

func (c *P2pClient) CreateVisitor(name string, localPort int) (string, int, error) {
	ip := "127.0.0.1"
	proxyConfigurer := &v1.XTCPProxyConfig{
		ProxyBaseConfig: v1.ProxyBaseConfig{
			Name: name,
			Type: "tcp",
			ProxyBackend: v1.ProxyBackend{
				LocalIP:   ip,
				LocalPort: localPort,
			},
		},
		Secretkey: "abcdefg",
	}
	proxyConfigurer.Complete("")

	c.pxyCfgs = append(c.pxyCfgs, proxyConfigurer)

	err := c.svr.ReloadConf(c.pxyCfgs, c.visitorCfgs)

	return ip, localPort, err
}

func handleTermSignal(svr *client.Service) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	svr.GracefulClose(500 * time.Millisecond)
}

func generate(serverAddr string, serverPort int) (*v1.ClientCommonConfig, []v1.ProxyConfigurer, []v1.VisitorConfigurer, bool, error) {

	cfg := &v1.ClientCommonConfig{
		ServerAddr: serverAddr,
		ServerPort: serverPort,
	}
	cfg.Complete()

	pxyCfgs := make([]v1.ProxyConfigurer, 0)
	visitorCfgs := make([]v1.VisitorConfigurer, 0)

	return cfg, pxyCfgs, visitorCfgs, false, nil
}
