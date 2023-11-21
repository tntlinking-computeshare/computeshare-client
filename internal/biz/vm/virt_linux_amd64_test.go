package vm

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func getVirtManager() *VirtManager {
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", "1",
		"service.name", "Name",
		"service.version", "Version",
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	manage, err := NewVirtManager(logger)
	if err != nil {
		panic(err)
	}
	return manage
}

func TestCreateVm(t *testing.T) {
	manage := getVirtManager()
	id, err := manage.Create("ubuntu1", "ubuntu-20.04", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC2mLWYddGeahdk6i3muy72XDbppnG4LIDhyj/rSuzLstdVLI7mF7efkwCZgyYcYRJoIjNI5mnb17o7/qVWdgGSiMnSgiPcw4r0Dp1pghWXBEog3o7pI3gicY6//Y4+liqypBEDmBSJnDsMJqVARzFV0rjJLhYSCbYk99LPB1ZLj0mDvIY/1SjRR9bfPuW9Ht6QjkS9DEWIdTrJ0dAaGwJkc+a5pCVzcopq4ycvBVLEnEq4xCrhbNx/LrpYxytA7WXg6kUcN+4Me63QVPxUExcn14qXr5uYxo+ePkoBCNdbqFsm0Z1rxrEX8oGDHvAfsoCpQr/OV8J5WwO7i/QIOyK7 mohaijiang110@163.com", "1", "1024")
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

func TestVirtManager_Status(t *testing.T) {
	manage := getVirtManager()

	status, err := manage.Status("ubuntu1")
	if err != nil {
		return
	}

	assert.NoError(t, err)
	fmt.Println(status)
}

func TestVirtManager_GetIp(t *testing.T) {
	manage := getVirtManager()
	ip, err := manage.GetIp("ubuntu1")
	if err != nil {
		return
	}

	assert.NoError(t, err)
	fmt.Println(ip)
}

func TestVirtManager_Reboot(t *testing.T) {
	manage := getVirtManager()

	err := manage.Reboot("ubuntu1")

	assert.NoError(t, err)
}

func TestVirtManager_Destroy(t *testing.T) {
	manage := getVirtManager()

	err := manage.Destroy("ubuntu1")
	assert.NoError(t, err)
}

func TestVirtManager_Init(t *testing.T) {
	manage := getVirtManager()

	manage.initBaseData()
}
