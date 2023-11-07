package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/mohaijiang/computeshare-client/internal/biz"
	"github.com/mohaijiang/computeshare-client/third_party/agent"
	queueTaskV1 "github.com/mohaijiang/computeshare-server/api/queue/v1"
	"time"
)

type CronJob struct {
	vmDockerService *VmDockerService
	agentService    *agent.AgentService
	p2pClient       *biz.P2pClient
	log             *log.Helper
}

func NewCronJob(vmDockerService *VmDockerService, agentService *agent.AgentService, p2pClient *biz.P2pClient, logger log.Logger) *CronJob {
	return &CronJob{
		vmDockerService: vmDockerService,
		agentService:    agentService,
		p2pClient:       p2pClient,
		log:             log.NewHelper(logger),
	}
}

func (c *CronJob) StartJob() {
	// 定时同步虚拟机的cpu和内存使用情况
	go c.syncInstanceStatus()

	go c.handlerQueueTask()
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
			c.vmDockerService.SyncServerVm()
			log.Info("结束同步虚拟机状态")
		}
	}
}

// handlerQueueTask 处理需要执行的命令
func (c *CronJob) handlerQueueTask() {
	// 创建一个定时触发的通道，每隔一秒发送一个时间事件
	ticker := time.Tick(5 * time.Second)

	// 使用 for 循环执行定时任务
	for {
		select {
		case <-ticker:
			// 在这里执行你的定时任务代码
			log.Debug("开始查询需要处理的任务")
			resp, err := c.agentService.GetQueueTask()
			if err != nil {
				return
			}
			task := resp.Data
			if task == nil {
				return
			}
			log.Infof("接收到需要处理的任务，任务id: %d, 任务类型： %s", task.Id, task.GetStatus())

			err = c.agentService.UpdateQueueTaskStatus(task.Id, queueTaskV1.TaskStatus_EXECUTING)
			if err != nil {
				log.Error("更新任务状态失败，跳过次任务： 异常内容：", err.Error())
				return
			}

			switch task.Cmd {
			case queueTaskV1.QueueCmd_VM_CREATE:
				{
				}
			case queueTaskV1.QueueCmd_VM_DELETE:
				{

				}
			case queueTaskV1.QueueCmd_VM_START:
				{

				}
			case queueTaskV1.QueueCmd_VM_SHUTDOWN:
				{
				}
			case queueTaskV1.QueueCmd_VM_RESTART:
				{

				}
			case queueTaskV1.QueueCmd_NAT_PROXY_CREATE:

				params := task.GetParams()

				var createParam queueTaskV1.NatProxyCreateVO
				err = json.Unmarshal([]byte(params), &createParam)

				if err == nil {
					_, _, err = c.p2pClient.CreateProxy(createParam.Name, "127.0.0.1", createParam.InstancePort, 6000)
				}

			case queueTaskV1.QueueCmd_NAT_PROXY_DELETE:
				var createParam queueTaskV1.NatProxyCreateVO
				err = json.Unmarshal([]byte(task.GetParams()), &createParam)
				if err == nil {
					err = c.p2pClient.DeleteProxy(createParam.Name)
				}
			case queueTaskV1.QueueCmd_NAT_VISITOR_CREATE:
				{

				}
			case queueTaskV1.QueueCmd_NAT_VISITOR_DELETE:
				{

				}
			default:
				log.Infof("无法确定执行任务命令，执行任务失败，任务id: %d", task.Id)
				err = fmt.Errorf("无法确定执行任务命令，执行任务失败，任务id: %d", task.Id)
			}

			if err != nil {
				_ = c.agentService.UpdateQueueTaskStatus(task.Id, queueTaskV1.TaskStatus_FAILED)
			} else {
				_ = c.agentService.UpdateQueueTaskStatus(task.Id, queueTaskV1.TaskStatus_EXECUTED)
			}
			log.Debug("结束任务处理")
		}
	}
}
