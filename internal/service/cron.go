package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type CronJob struct {
	vmService *VmService
	log       *log.Helper
}

func NewCronJob(vmService *VmService, logger log.Logger) *CronJob {
	return &CronJob{
		vmService: vmService,
		log:       log.NewHelper(logger),
	}
}

func (c *CronJob) StartJob() {
	// 定时同步虚拟机的cpu和内存使用情况
	go c.syncInstanceStatus()
}

func (c *CronJob) syncInstanceStatus() {
	// 创建一个定时触发的通道，每隔一秒发送一个时间事件
	ticker := time.Tick(1 * time.Minute)

	// 使用 for 循环执行定时任务
	for {
		select {
		case <-ticker:
			// 在这里执行你的定时任务代码
			log.Info("开始同步虚拟机状态")
			c.vmService.SyncServerVm()
			log.Info("结束同步虚拟机状态")
		}
	}
}
