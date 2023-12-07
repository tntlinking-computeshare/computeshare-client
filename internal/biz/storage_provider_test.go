package biz

import (
	queueTaskV1 "github.com/mohaijiang/computeshare-server/api/queue/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunSeaweedFSContainer(t *testing.T) {
	param := &queueTaskV1.StorageSetupTaskParamVO{
		MasterServer: "computeshare.newtouch.com:9333",
		PublicIp:     "computeshare.newtouch.com",
		PublicPort:   41016,
		GrpcPort:     41017,
	}
	err := runSeaweedFSContainer(param)

	assert.NoError(t, err)
}

func TestStopSeaweedFSContainer(t *testing.T) {
	err := stopSeaweedFSContainer()
	assert.NoError(t, err)

}
