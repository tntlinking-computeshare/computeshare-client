package biz

import (
	"fmt"
	"testing"
)

func TestSaveConfig(t *testing.T) {
	cfg := FrpClientConfigure{
		ServerPort: 7000,
		ServerAddr: "61.172.179.6",
		Visitors: []Visitor{
			{
				Name:           "p2p_ssh_visitor",
				Type:           "xtcp",
				ServerName:     "p2p_ssh",
				SecretKey:      "abcdefg",
				BindAddr:       "127.0.0.1",
				BindPort:       6000,
				KeepTunnelOpen: false,
			},
		},
	}

	cfg.Save("/tmp/frpc.toml")
}

func TestSaveConfigure2(t *testing.T) {
	configure := &FrpClientConfigure{
		ServerAddr: "61.172.179.73",
		ServerPort: 7000,
	}

	for i := 9001; i <= 9050; i++ {
		configure.Proxies = append(configure.Proxies, Proxy{
			Name:       fmt.Sprintf("yiliao_%d", i),
			Type:       "tcp",
			LocalIP:    "192.168.122.72",
			LocalPort:  i,
			RemotePort: i,
		})
	}

	configure.Println()

}
