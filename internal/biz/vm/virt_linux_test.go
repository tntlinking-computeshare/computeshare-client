package vm

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/tntlinking-computeshare/computeshare-client/internal/conf"
	"github.com/tntlinking-computeshare/computeshare-client/third_party/agent"
	queueTaskV1 "github.com/tntlinking-computeshare/computeshare-server/api/server/queue/v1"
	"os"
	"testing"
)

func getVirtManager() IVirtManager {
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", "1",
		"service.name", "Name",
		"service.version", "Version",
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	data := &conf.Data{}
	httpClient, _, err := agent.NewHttpConnection(data)
	manage, err := NewVirtManager(logger, cli, data, httpClient)
	if err != nil {
		panic(err)
	}
	return manage
}

func TestCreateVm(t *testing.T) {
	manage := getVirtManager()
	param := &queueTaskV1.ComputeInstanceTaskParamVO{
		Id:            "myInstanceId",
		InstanceId:    "myInstanceId",
		Name:          "ubuntu1",
		Image:         "ubuntu:20.04",
		Password:      "Abcd1234",
		PublicKey:     "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC2mLWYddGeahdk6i3muy72XDbppnG4LIDhyj/rSuzLstdVLI7mF7efkwCZgyYcYRJoIjNI5mnb17o7/qVWdgGSiMnSgiPcw4r0Dp1pghWXBEog3o7pI3gicY6//Y4+liqypBEDmBSJnDsMJqVARzFV0rjJLhYSCbYk99LPB1ZLj0mDvIY/1SjRR9bfPuW9Ht6QjkS9DEWIdTrJ0dAaGwJkc+a5pCVzcopq4ycvBVLEnEq4xCrhbNx/LrpYxytA7WXg6kUcN+4Me63QVPxUExcn14qXr5uYxo+ePkoBCNdbqFsm0Z1rxrEX8oGDHvAfsoCpQr/OV8J5WwO7i/QIOyK7 mohaijiang110@163.com",
		Cpu:           1,
		Memory:        2,
		DockerCompose: "dmVyc2lvbjogIjMiCnNlcnZpY2VzOgogIHdlYjoKICAgIGltYWdlOiBuZ2lueDpsYXRlc3QKICAgIHBvcnRzOgogICAgICAtICI4MDo4MCI=",
	}
	id, err := manage.Create(param)
	if err != nil {
		panic(err)
	}

	fmt.Println(id)
}

func TestVirtManager_Shutdown(t *testing.T) {
	manage := getVirtManager()

	err := manage.Shutdown("ubuntu1")

	assert.NoError(t, err)
}

func TestVirtManager_Start(t *testing.T) {
	manage := getVirtManager()

	err := manage.Start("ubuntu1")

	assert.NoError(t, err)
}

func TestStatus(t *testing.T) {
	manage := getVirtManager()
	ip, err := manage.GetIp("3043870d-c25d-4733-84f1-2a8a6a0a6ada")
	if err != nil {
		panic(ip)
	}
	fmt.Println(ip)
}

func TestVirtManager_Destroy(t *testing.T) {
	manage := getVirtManager()

	err := manage.Destroy("ubuntu1")
	assert.NoError(t, err)
}

func TestTemplate(t *testing.T) {
	manage := getVirtManager()
	systemInfo, err := manage.GetSystemInfo()
	assert.NoError(t, err)
	fmt.Println(systemInfo)
}
