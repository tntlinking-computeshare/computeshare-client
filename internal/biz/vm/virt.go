package vm

import (
	queueTaskV1 "github.com/mohaijiang/computeshare-server/api/queue/v1"
)

type IVirtManager interface {
	Create(param *queueTaskV1.ComputeInstanceTaskParamVO) (string, error)
	Destroy(name string) error
	Start(name string) error
	Shutdown(name string) error
	Reboot(name string) error
	VncOpen(name string) (int32, error)
	GetIp(name string) (string, error)
}
