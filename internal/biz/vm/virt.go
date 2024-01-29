package vm

import (
	"github.com/libvirt/libvirt-go"
	queueTaskV1 "github.com/mohaijiang/computeshare-server/api/queue/v1"
)

type IVirtManager interface {
	Create(param *queueTaskV1.ComputeInstanceTaskParamVO) (string, error)
	Destroy(name string) error
	Start(name string) error
	Shutdown(name string) error
	Reboot(name string) error
	VncOpen(name string, vncPort int32) error
	GetIp(name string) (string, error)
	GetVncWebsocketIP(name string) (string, error)
	GetVncWebsocketPort(name string) int32
	ReCreate(name string, param *queueTaskV1.ComputeInstanceTaskParamVO) error
	Status(name string) (libvirt.DomainState, error)
}
