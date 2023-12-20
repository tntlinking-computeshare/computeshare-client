package agent

import (
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMac(t *testing.T) {

	ip, mac, err := getLocalIPAndMacAddress()
	assert.NoError(t, err)
	fmt.Println("ip: ", ip)
	fmt.Println("mac:", mac)

	info, _ := mem.VirtualMemory()
	fmt.Println(info.Total / 1024 / 1024)
}
