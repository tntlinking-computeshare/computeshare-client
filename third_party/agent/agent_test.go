package agent

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMac(t *testing.T) {

	address, err := getLocalMacAddress()
	assert.NoError(t, err)
	fmt.Println(address)
}
