package p2p

import (
	"encoding/base64"
	"fmt"
	"github.com/google/wire"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/mohaijiang/computeshare-client/internal/conf"
	goipfsp2p "github.com/mohaijiang/go-ipfs-p2p"
	"os"
	"path"
)

var ProviderSet = wire.NewSet(NewP2pClient)

func NewP2pClient(c *conf.Server) (*goipfsp2p.P2pClient, error) {
	basePath := os.Getenv("IPFS_PATH")

	if basePath == "" {
		home, _ := os.UserHomeDir()
		basePath = path.Join(home, ".ipfs")
	}

	_ = os.MkdirAll(basePath, 0644)

	privateKeyPath := path.Join(basePath, "privateKey")
	privateKeyByte, err := os.ReadFile(privateKeyPath)
	privateKey := string(privateKeyByte)

	if !CheckPrivateKeyAvaliable(privateKey) || err != nil {
		priv, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
		skbytes, err := crypto.MarshalPrivateKey(priv)
		if err != nil {
			panic(err)
		}
		privateKey = base64.StdEncoding.EncodeToString(skbytes)
		_ = os.WriteFile(privateKeyPath, []byte(privateKey), 0644)
	}

	return goipfsp2p.NewP2pClient(int(c.P2P.Port), privateKey, c.P2P.SwarmKey, c.P2P.GetBootstraps())
}

func CheckPrivateKeyAvaliable(privateKey string) bool {
	//检查privateKey 是否有效
	skbytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		fmt.Println(err)
		return false

	}
	_, err = crypto.UnmarshalPrivateKey(skbytes)
	return err == nil
}
