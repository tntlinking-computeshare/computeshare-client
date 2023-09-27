package agent

import (
	"context"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	shell "github.com/ipfs/go-ipfs-api"
	agentv1 "github.com/mohaijiang/computeshare-server/api/agent/v1"
	"time"
)

type AgentService struct {
	client    agentv1.AgentHTTPClient
	ipfsShell *shell.Shell
	id        string
}

func NewAgentService(conn *transhttp.Client, ipfsShell *shell.Shell) *AgentService {

	//client := pb.New(conn)
	client := agentv1.NewAgentHTTPClient(conn)
	return &AgentService{
		client:    client,
		ipfsShell: ipfsShell,
	}
}

func (s *AgentService) Register() error {
	peerId := ""
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
