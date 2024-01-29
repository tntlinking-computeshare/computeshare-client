package agent

import (
	"context"
	"errors"
	"fmt"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	agentv1 "github.com/mohaijiang/computeshare-server/api/agent/v1"
	"github.com/mohaijiang/computeshare-server/api/compute/v1"
	queueTaskV1 "github.com/mohaijiang/computeshare-server/api/queue/v1"
	"github.com/shirou/gopsutil/mem"
	"net"
	"os"
	"runtime"
	"time"
)

type AgentService struct {
	client          agentv1.AgentHTTPClient
	queueTaskClient queueTaskV1.QueueTaskHTTPClient
	id              string
}

func NewAgentService(conn *transhttp.Client) *AgentService {

	client := agentv1.NewAgentHTTPClient(conn)
	queueTaskClient := queueTaskV1.NewQueueTaskHTTPClient(conn)
	return &AgentService{
		client:          client,
		queueTaskClient: queueTaskClient,
	}
}

func (s *AgentService) Register() error {
	ip, mac, err := getLocalIPAndMacAddress()
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	info, _ := mem.VirtualMemory()
	totalMemory := info.Total / (1024 * 1024 * 1024) // 转换为GB
	used := info.Used / (1024 * 1024 * 1024)

	res, err := s.client.CreateAgent(ctx, &agentv1.CreateAgentRequest{
		Mac:            mac,
		Hostname:       hostname,
		TotalCpu:       int32(runtime.NumCPU()),
		TotalMemory:    int32(totalMemory),
		OccupiedMemory: int32(used),
		Ip:             ip,
	})

	if err != nil {
		return err
	}

	s.id = res.Data.Id

	return nil
}

func (s *AgentService) UnRegister() error {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)

	_, err := s.client.DeleteAgent(ctx, &agentv1.DeleteAgentRequest{
		Id: s.id,
	})

	return err
}

func (s *AgentService) ListInstances() (*v1.ListInstanceReply, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	return s.client.ListAgentInstance(ctx, &agentv1.ListAgentInstanceReq{
		AgentId: s.id,
	})
}

func (s *AgentService) ReportContainerStatus(instance *v1.Instance) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	_, err := s.client.ReportInstanceStatus(ctx, instance)
	return err
}

func (s *AgentService) GetMac() string {
	_, mac, err := getLocalIPAndMacAddress()
	if err != nil {
		return ""
	}
	return mac
}

func (s *AgentService) GetQueueTask() (*queueTaskV1.QueueTaskGetResponse, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	return s.queueTaskClient.GetAgentTask(ctx, &queueTaskV1.QueueTaskGetRequest{
		Id: s.id,
	})
}

func (s *AgentService) UpdateQueueTaskStatus(taskId string, status queueTaskV1.TaskStatus) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)

	_, err := s.queueTaskClient.UpdateAgentTask(ctx, &queueTaskV1.QueueTaskUpdateRequest{
		Id:      taskId,
		AgentId: s.id,
		Status:  status,
	})
	return err
}

func getLocalIPAndMacAddress() (string, string, error) {
	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", "", err
	}
	for _, iface := range interfaces {

		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 获取网络接口的IP地址
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("无法获取网络接口地址:", err)
			continue
		}

		// 输出第一个非环回地址
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				fmt.Println("无法解析IP地址:", err)
				continue
			}

			if ip.To4() != nil {
				fmt.Printf("当前IP地址: %s\n", ip)
				return ip.String(), iface.HardwareAddr.String(), nil
			}
		}
	}
	return "", "", errors.New("cannot get mac address")
}
