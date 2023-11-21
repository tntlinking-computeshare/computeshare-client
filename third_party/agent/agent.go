package agent

import (
	"context"
	"errors"
	"fmt"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	agentv1 "github.com/mohaijiang/computeshare-server/api/agent/v1"
	"github.com/mohaijiang/computeshare-server/api/compute/v1"
	queueTaskV1 "github.com/mohaijiang/computeshare-server/api/queue/v1"
	"net"
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
	peerId, err := getLocalMacAddress()
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	res, err := s.client.CreateAgent(ctx, &agentv1.CreateAgentRequest{
		Name: peerId,
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
		PeerId: s.GetPeerId(),
	})
}

func (s *AgentService) ReportContainerStatus(instance *v1.Instance) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	_, err := s.client.ReportInstanceStatus(ctx, instance)
	return err
}

func (s *AgentService) GetPeerId() string {
	peerId, err := getLocalMacAddress()
	if err != nil {
		return ""
	}
	return peerId
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

func getLocalMacAddress() (string, error) {
	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, inter := range interfaces {
		mac := inter.HardwareAddr //获取本机MAC地址
		if mac.String() != "" {
			fmt.Println(inter.Name)
			fmt.Println("MAC = ", mac)
			return mac.String(), nil
		}
	}
	return "", errors.New("cannot get mac address")
}
