package biz

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/naoina/toml"
	"github.com/samber/lo"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

type FrpClientConfigure struct {
	ServerAddr string `json:"serverAddr" toml:"serverAddr"`
	ServerPort int    `json:"serverPort" toml:"serverPort"`

	Proxies  []Proxy   `json:"proxies" toml:"proxies,omitempty"`
	Visitors []Visitor `json:"visitors" toml:"visitors,omitempty"`
}

func (c *FrpClientConfigure) Save(path string) {
	marshal, err := toml.Marshal(c)
	if err != nil {
		return
	}

	err = os.WriteFile(path, marshal, 0644)
	if err != nil {
		return
	}
}

func LoadFrpClientConfigure(path string) (*FrpClientConfigure, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg FrpClientConfigure

	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, err
}

type Proxy struct {
	Name       string `json:"name" toml:"name"`
	Type       string `json:"type" toml:"type"`
	LocalIP    string `json:"localIP" toml:"localIP"`
	LocalPort  int    `json:"localPort" toml:"localPort"`
	RemotePort int    `json:"remotePort" toml:"remotePort"`
}

type Visitor struct {
	Name           string `json:"name" toml:"name"`
	Type           string `json:"type" toml:"type"`
	ServerName     string `json:"serverName" toml:"serverName"`
	SecretKey      string `json:"secretKey" toml:"SecretKey"`
	BindAddr       string `json:"bindAddr" toml:"BindAddr"`
	BindPort       int    `json:"bindPort" toml:"bindPort"`
	KeepTunnelOpen bool   `json:"keepTunnelOpen" toml:"keepTunnelOpen"`
}

type P2pClient struct {
	svr         *client.Service
	pxyCfgs     []v1.ProxyConfigurer
	visitorCfgs []v1.VisitorConfigurer
	gatewayIp   string
	gatewayPort int

	configPath string
}

func NewP2pClient() *P2pClient {
	dir, _ := os.UserHomeDir()

	c := &P2pClient{
		configPath: path.Join(dir, ".frpc.toml"),
	}

	if configure, err := LoadFrpClientConfigure(c.configPath); err == nil {
		_ = c.Start(configure.ServerAddr, int32(configure.ServerPort))
	}

	return c
}

func (c *P2pClient) IsStart() bool {
	return c.svr != nil
}

func (c *P2pClient) Start(gatewayIp string, gatewayPort int32) error {

	if c.IsStart() {
		return nil
	}

	configure, err := LoadFrpClientConfigure(c.configPath)
	if err != nil {
		configure = &FrpClientConfigure{
			ServerAddr: gatewayIp,
			ServerPort: int(gatewayPort),
		}

		configure.Save(c.configPath)
	}

	cfg, pxyCfgs, visitorCfgs, _, err := config.LoadClientConfig(c.configPath)
	if err != nil {
		return err
	}

	svr, err := client.NewService(cfg, pxyCfgs, visitorCfgs, c.configPath)
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

func (c *P2pClient) CreateProxy(name string, localIp string, localPort, remotePort int32, protocol string) (string, int, error) {

	if c.svr == nil {
		return "", 0, errors.New("no initialization")
	}

	configure, err := LoadFrpClientConfigure(c.configPath)
	if err != nil {
		return "", 0, err
	}

	proxyType := "tcp"
	if protocol == "UDP" {
		proxyType = "udp"
	}

	configure.Proxies = append(configure.Proxies, Proxy{
		Name:       name,
		Type:       proxyType,
		LocalIP:    localIp,
		LocalPort:  int(localPort),
		RemotePort: int(remotePort),
	})

	configure.Save(c.configPath)

	_, pxyCfgs, _, _, err := config.LoadClientConfig(c.configPath)
	if err != nil {
		return "", 0, err
	}

	c.pxyCfgs = pxyCfgs

	err = c.svr.ReloadConf(c.pxyCfgs, c.visitorCfgs)

	return c.gatewayIp, int(remotePort), err
}

func (c *P2pClient) EditProxy(name string, localIp string, localPort int32, remotePort int32, protocol string) error {
	configure, err := LoadFrpClientConfigure(c.configPath)
	if err != nil {
		return err
	}
	proxy, ok := lo.Find(configure.Proxies, func(item Proxy) bool {
		return item.Name == name
	})

	proxyType := "tcp"
	if protocol == "UDP" {
		proxyType = "udp"
	}
	if !ok {
		return nil
	}

	proxy.Type = proxyType
	proxy.LocalIP = localIp
	proxy.LocalPort = int(localPort)
	proxy.RemotePort = int(remotePort)

	// 删除原端口映射
	configure.Proxies = lo.FilterMap(configure.Proxies, func(item Proxy, _ int) (Proxy, bool) {
		return item, item.Name != name
	})
	// 重新保存端口映射
	configure.Proxies = append(configure.Proxies, proxy)
	configure.Save(c.configPath)

	return c.svr.ReloadConf(c.pxyCfgs, c.visitorCfgs)
}

func (c *P2pClient) DeleteProxy(name string) error {
	if c.svr == nil {
		return errors.New("no initialization")
	}

	c.pxyCfgs = lo.FilterMap(c.pxyCfgs, func(item v1.ProxyConfigurer, _ int) (v1.ProxyConfigurer, bool) {
		fmt.Println(item.GetBaseConfig().Name, "==", name)
		return item, item.GetBaseConfig().Name != name
	})

	configure, err := LoadFrpClientConfigure(c.configPath)
	if err != nil {
		return err
	}
	configure.Proxies = lo.FilterMap(configure.Proxies, func(item Proxy, _ int) (Proxy, bool) {
		return item, item.Name != name
	})
	configure.Save(c.configPath)

	return c.svr.ReloadConf(c.pxyCfgs, c.visitorCfgs)
}

func (c *P2pClient) CreateVisitor(name string, localPort int) (string, int, error) {
	if c.svr == nil {
		return "", 0, errors.New("no initialization")
	}

	configure, err := LoadFrpClientConfigure(c.configPath)
	if err != nil {
		return "", 0, err
	}

	configure.Visitors = append(configure.Visitors, Visitor{
		Name:       name,
		Type:       "tcp",
		ServerName: "",
		SecretKey:  "abcdefg",
		BindAddr:   "127.0.0.1",
		BindPort:   localPort,
	})

	configure.Save(c.configPath)

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

	err = c.svr.ReloadConf(c.pxyCfgs, c.visitorCfgs)

	return ip, localPort, err
}

func handleTermSignal(svr *client.Service) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	svr.GracefulClose(500 * time.Millisecond)
}

//func generate(serverAddr string, serverPort int32) (*v1.ClientCommonConfig, []v1.ProxyConfigurer, []v1.VisitorConfigurer, bool, error) {
//
//	cfg := &v1.ClientCommonConfig{
//		ServerAddr: serverAddr,
//		ServerPort: int(serverPort),
//	}
//	cfg.Complete()
//
//	pxyCfgs := make([]v1.ProxyConfigurer, 0)
//	visitorCfgs := make([]v1.VisitorConfigurer, 0)
//
//	return cfg, pxyCfgs, visitorCfgs, false, nil
//}
