package agent

import (
	"context"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	agentv1 "github.com/mohaijiang/computeshare-server/api/agent/v1"
	go_ipfs_p2p "github.com/mohaijiang/go-ipfs-p2p"
	"time"
)

type AgentService struct {
	client    agentv1.AgentHTTPClient
	id        string
	p2pClient *go_ipfs_p2p.P2pClient
}

func NewAgentService(conn *transhttp.Client, p2pClient *go_ipfs_p2p.P2pClient) *AgentService {

	//client := pb.New(conn)
	client := agentv1.NewAgentHTTPClient(conn)
	return &AgentService{
		client:    client,
		p2pClient: p2pClient,
	}
}

func (s *AgentService) Register() error {
	peerId := s.p2pClient.Host.ID().String()
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
