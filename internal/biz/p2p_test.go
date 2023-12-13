package biz

import (
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
