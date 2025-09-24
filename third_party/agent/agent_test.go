package agent

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMac(t *testing.T) {

	ip, mac, err := getLocalIPAndMacAddress()
	assert.NoError(t, err)
	fmt.Println("ip: ", ip)
	fmt.Println("mac:", mac)
}
