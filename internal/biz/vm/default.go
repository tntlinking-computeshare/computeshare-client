// //go:build default
package vm

//
//import (
//	"errors"
//	"github.com/go-kratos/kratos/v2/log"
//	queueTaskV1 "github.com/mohaijiang/computeshare-server/api/queue/v1"
//)
//
//func NewVirtManager(logger log.Logger) (IVirtManager, error) {
//	return NewUnsupportVirtManager(), nil
//}
//
//type UnsupportVirtManager struct {
//}
//
//func (u UnsupportVirtManager) Create(param *queueTaskV1.ComputeInstanceTaskParamVO) (string, error) {
//	//TODO implement me
//	return "", errors.New("unsupport")
//}
//
//func (u UnsupportVirtManager) Destroy(name string) error {
//	//TODO implement me
//	return errors.New("unsupport")
//}
//
//func (u UnsupportVirtManager) Start(name string) error {
//	//TODO implement me
//	return errors.New("unsupport")
//}
//
//func (u UnsupportVirtManager) Shutdown(name string) error {
//	//TODO implement me
//	return errors.New("unsupport")
//}
//
//func (u UnsupportVirtManager) Reboot(name string) error {
//	//TODO implement me
//	return errors.New("unsupport")
//}
//
//func (u UnsupportVirtManager) VncOpen(name string) (int32, error) {
//	//TODO implement me
//	return 0, errors.New("unsupport")
//}
//
//func (u UnsupportVirtManager) GetIp(name string) (string, error) {
//	//TODO implement me
//	return "", errors.New("unsupport")
//}
//
//func (u UnsupportVirtManager) GetVncWebsocketIP(name string) (string, error) {
//	return "", errors.New("unsupport")
//}
//func (u UnsupportVirtManager) GetVncWebsocketPort(name string) int32 {
//	return 0
//}
//
//func NewUnsupportVirtManager() IVirtManager {
//	return &UnsupportVirtManager{}
//}
