package service

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/libvirt/libvirt-go"
	"github.com/mohaijiang/computeshare-client/internal/biz"
	"github.com/mohaijiang/computeshare-client/internal/biz/vm"
	"github.com/mohaijiang/computeshare-client/third_party/agent"
	"github.com/mohaijiang/computeshare-server/api/compute"
	queueTaskV1 "github.com/mohaijiang/computeshare-server/api/queue/v1"
	"time"
)

type CronJob struct {
	agentService    *agent.AgentService
	p2pClient       *biz.P2pClient
	virtManager     vm.IVirtManager
	storageProvider *biz.StorageProvider
	log             *log.Helper
}

func NewCronJob(
	agentService *agent.AgentService,
	p2pClient *biz.P2pClient,
	virtManager vm.IVirtManager,
	storageProvider *biz.StorageProvider,
	logger log.Logger) *CronJob {
	return &CronJob{
		agentService:    agentService,
		p2pClient:       p2pClient,
		virtManager:     virtManager,
		storageProvider: storageProvider,
		log:             log.NewHelper(logger),
	}
}

func (c *CronJob) StartJob() {

	// 先同步虚拟机状态
	c.SyncComputeInstanceStatus()
	// 同步虚拟机任务队列
	go c.handlerQueueTask()
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
			taskResp, err := c.agentService.GetQueueTask()
			if err != nil {
				continue
			}
			c.DoTask(taskResp)
		}
	}
}

func (c *CronJob) DoTask(taskResp *queueTaskV1.QueueTaskGetResponse) {
	task := taskResp.Data
	if task == nil {
		return
	}
	log.Infof("接收到需要处理的任务，任务id: %d, 任务类型： %s", task.Id, task.GetCmd())

	err := c.agentService.UpdateQueueTaskStatus(task.Id, queueTaskV1.TaskStatus_EXECUTING)
	if err != nil {
		log.Error("更新任务状态失败，跳过次任务： 异常内容：", err.Error())
		return
	}

	params, jsonErr := task.GetTaskParam()
	if jsonErr != nil {
		log.Error("解析任务参数失败, 任务参数：", task.GetParams())
		_ = c.agentService.UpdateQueueTaskStatus(task.Id, queueTaskV1.TaskStatus_FAILED)
		return
	}
	switch task.Cmd {
	case queueTaskV1.TaskCmd_VM_CREATE:
		createParam, ok := params.(*queueTaskV1.ComputeInstanceTaskParamVO)
		if ok {
			_, err = c.virtManager.Create(createParam)

			if !c.p2pClient.IsStart() {
				err = c.p2pClient.Start(createParam.GatewayIp, createParam.GatewayPort)
				if err != nil {
					break
				}
			}
			var localIP string
			localIP, err = c.virtManager.GetVncWebsocketIP(createParam.InstanceId)
			if err != nil {
				fmt.Println(err)
				break
			}

			localPort := c.virtManager.GetVncWebsocketPort(createParam.InstanceId)

			_, _, err = c.p2pClient.CreateProxy(fmt.Sprintf("vnc_%s", createParam.InstanceId), localIP, localPort, createParam.VncConnectPort, "tcp")

		}
	case queueTaskV1.TaskCmd_VM_DELETE:
		createParam, ok := params.(*queueTaskV1.ComputeInstanceTaskParamVO)
		if ok {
			err = c.virtManager.Destroy(createParam.Id)
			if !c.p2pClient.IsStart() {
				err = c.p2pClient.Start(createParam.GatewayIp, createParam.GatewayPort)
				if err != nil {
					fmt.Println(err)
					break
				}
			}
			err = c.p2pClient.DeleteProxy(fmt.Sprintf("vnc_%s", createParam.InstanceId))
		}
	case queueTaskV1.TaskCmd_VM_START:
		createParam, ok := params.(*queueTaskV1.ComputeInstanceTaskParamVO)
		if ok {
			err = c.virtManager.Start(createParam.Id)
		}
	case queueTaskV1.TaskCmd_VM_SHUTDOWN:
		createParam, ok := params.(*queueTaskV1.ComputeInstanceTaskParamVO)
		if ok {
			err = c.virtManager.Shutdown(createParam.Id)
		}
	case queueTaskV1.TaskCmd_VM_RESTART:
		createParam, ok := params.(*queueTaskV1.ComputeInstanceTaskParamVO)
		if ok {
			err = c.virtManager.Reboot(createParam.Id)
		}
	case queueTaskV1.TaskCmd_VM_RECREATE:
		createParam, ok := params.(*queueTaskV1.ComputeInstanceTaskParamVO)
		if ok {
			err = c.virtManager.ReCreate(createParam.InstanceId, createParam)
			if err != nil {
				return
			}
		}
	case queueTaskV1.TaskCmd_NAT_PROXY_CREATE:

		createParam, ok := params.(*queueTaskV1.NatNetworkMappingTaskParamVO)
		if ok {
			if !c.p2pClient.IsStart() {
				err = c.p2pClient.Start(createParam.GatewayIp, createParam.GatewayPort)
				if err != nil {
					_ = c.agentService.UpdateQueueTaskStatus(task.Id, queueTaskV1.TaskStatus_FAILED)
					fmt.Println(err)
					break
				}
			}

			instanceName := createParam.InstanceId

			var ip string
			if instanceName != "" {
				ip, err = c.virtManager.GetIp(instanceName)
			} else {
				ip = "127.0.0.1"
			}
			if err == nil {
				_, _, err = c.p2pClient.CreateProxy(createParam.Name, ip, createParam.InstancePort, createParam.RemotePort, createParam.Protocol)
			}
		}

	case queueTaskV1.TaskCmd_NAT_PROXY_DELETE:
		params, jsonErr := task.GetTaskParam()
		if jsonErr == nil {
			createParam, ok := params.(*queueTaskV1.NatNetworkMappingTaskParamVO)
			if ok {
				err = c.p2pClient.DeleteProxy(createParam.Name)
			}
		}
	case queueTaskV1.TaskCmd_NAT_PROXY_EDIT:
		params, jsonErr := task.GetTaskParam()
		if jsonErr == nil {
			createParam, ok := params.(*queueTaskV1.NatNetworkMappingTaskParamVO)
			var ip string
			if createParam.InstanceId != "" {
				ip, err = c.virtManager.GetIp(createParam.InstanceId)
			} else {
				ip = "127.0.0.1"
			}
			if ok {
				err = c.p2pClient.EditProxy(createParam.Name, ip, createParam.InstancePort, createParam.RemotePort, createParam.Protocol)
			}
		}
	case queueTaskV1.TaskCmd_NAT_VISITOR_CREATE:
		{

		}
	case queueTaskV1.TaskCmd_NAT_VISITOR_DELETE:
		{

		}

	case queueTaskV1.TaskCmd_STORAGE_CREATE:
		createParam, ok := params.(*queueTaskV1.StorageSetupTaskParamVO)
		if ok {
			if !c.storageProvider.Status() {
				err = c.storageProvider.Start(createParam)
			}
		}
	case queueTaskV1.TaskCmd_STORAGE_DELETE:
		err = c.storageProvider.Stop()

	default:
		log.Infof("无法确定执行任务命令，执行任务失败，任务id: %d", task.Id)
		err = fmt.Errorf("无法确定执行任务命令，执行任务失败，任务id: %d", task.Id)
	}

	if err != nil {
		log.Errorf("执行任务命令失败：，失败原因：%s ", err.Error())
		_ = c.agentService.UpdateQueueTaskStatus(task.Id, queueTaskV1.TaskStatus_FAILED)
	} else {
		_ = c.agentService.UpdateQueueTaskStatus(task.Id, queueTaskV1.TaskStatus_EXECUTED)
	}
	log.Debug("结束任务处理")
}

func (c *CronJob) SyncComputeInstanceStatus() {
	instances, err := c.agentService.ListInstances()
	if err != nil {
		c.log.Errorf("sync computeInstance fail: %v", err)
		return
	}

	for _, instance := range instances.Data {

		instanceId := instance.Id
		// 判断本地有无此 instance

		state, err := c.virtManager.Status(instanceId)

		if err != nil {
			c.log.Errorf("cannot get instance state")
			continue
		}

		status := compute.InstanceStatus(instance.GetStatus())

		if status == compute.InstanceStatusRunning && (state == libvirt.DOMAIN_SHUTOFF || state == libvirt.DOMAIN_SHUTDOWN) {
			_ = c.virtManager.Start(instanceId)
			continue
		}

		if status == compute.InstanceStatusClosed && state == libvirt.DOMAIN_RUNNING {
			_ = c.virtManager.Shutdown(instanceId)
			continue
		}

	}
}
