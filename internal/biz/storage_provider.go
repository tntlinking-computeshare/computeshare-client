package biz

import (
	"fmt"
	queueTaskV1 "github.com/mohaijiang/computeshare-server/api/queue/v1"
)

func NewStorageProvider() *StorageProvider {
	return &StorageProvider{
		status: false,
	}
}

type StorageProvider struct {
	status bool
}

func (sp *StorageProvider) Status() bool {
	return sp.status
}

func (sp *StorageProvider) Start(param *queueTaskV1.StorageSetupTaskParamVO) error {
	sp.status = true
	fmt.Println(param.PublicPort)

	return nil
}

func (sp *StorageProvider) Stop() error {
	sp.status = false
	return nil
}
