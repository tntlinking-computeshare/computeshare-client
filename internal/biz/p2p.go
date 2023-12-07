package biz

import (
	"context"
	"errors"
	"github.com/fatedier/frp/client"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/samber/lo"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type P2pClient struct {
	svr         *client.Service
	pxyCfgs     []v1.ProxyConfigurer
	visitorCfgs []v1.VisitorConfigurer
	gatewayIp   string
	gatewayPort int
}

func NewP2pClient() *P2pClient {
	return &P2pClient{}
}

func (c *P2pClient) IsStart() bool {
	return c.svr != nil
}

func (c *P2pClient) Start(gatewayIp string, gatewayPort int32) error {

	if c.IsStart() {
		return nil
	}

	cfg, pxyCfgs, visitorCfgs, _, err := generate(gatewayIp, gatewayPort)
	if err != nil {
		return err
	}

	svr, err := client.NewService(cfg, pxyCfgs, visitorCfgs, "")
	if err != nil {
		return err
	}
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
	c.svr = svr
	c.pxyCfgs = pxyCfgs
	c.visitorCfgs = visitorCfgs

	return nil
}

func (c *P2pClient) CreateProxy(name string, localIp string, localPort, remotePort int32) (string, int, error) {

	if c.svr == nil {
		return "", 0, errors.New("no initialization")
	}

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

	return c.gatewayIp, int(remotePort), err
}

func (c *P2pClient) DeleteProxy(name string) error {
	if c.svr == nil {
		return errors.New("no initialization")
	}

	c.pxyCfgs = lo.FilterMap(c.pxyCfgs, func(item v1.ProxyConfigurer, _ int) (v1.ProxyConfigurer, bool) {
		return item, item.GetBaseConfig().Name != name
	})

	return c.svr.ReloadConf(c.pxyCfgs, c.visitorCfgs)
}

func (c *P2pClient) CreateVisitor(name string, localPort int) (string, int, error) {
	if c.svr == nil {
		return "", 0, errors.New("no initialization")
	}
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

func generate(serverAddr string, serverPort int32) (*v1.ClientCommonConfig, []v1.ProxyConfigurer, []v1.VisitorConfigurer, bool, error) {

	cfg := &v1.ClientCommonConfig{
		ServerAddr: serverAddr,
		ServerPort: int(serverPort),
	}
	cfg.Complete()

	pxyCfgs := make([]v1.ProxyConfigurer, 0)
	visitorCfgs := make([]v1.VisitorConfigurer, 0)

	return cfg, pxyCfgs, visitorCfgs, false, nil
}
