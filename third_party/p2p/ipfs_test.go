package p2p

import (
	"fmt"
	"github.com/tj/assert"
	"testing"
)

func TestIpfsDaemon(t *testing.T) {
	node, err := RunDaemon()
	if err != nil {
		assert.NoError(t, err)
	}

	PeerId := node.Identity.String()

	fmt.Println(PeerId)

	localAddress := fmt.Sprintf("/ip4/%s/tcp/%d/p2p/%s", "127.0.0.1", 4001, node.Identity.String())
	fmt.Println(localAddress)

	select {}
}
