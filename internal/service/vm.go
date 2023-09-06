package service

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/lo"

	pb "computeshare-client/api/compute/v1"
)

type VmService struct {
	pb.UnimplementedVmServer

	cli *client.Client
	log *log.Helper
}

func NewVmService(client *client.Client, logger log.Logger) *VmService {
	return &VmService{
		cli: client,
		log: log.NewHelper(logger),
	}
}

func (s *VmService) CreateVm(ctx context.Context, req *pb.CreateVmRequest) (*pb.GetVmReply, error) {
	ctx2 := context.Background()
	out, err := s.cli.ImagePull(ctx2, req.Image, types.ImagePullOptions{})
	s.log.Info(out)
	if err != nil {
		return nil, err
	}

	port, err := nat.NewPort("tcp", req.GetPort())
	resp, err := s.cli.ContainerCreate(ctx2, &container.Config{
		Image: req.Image,
		ExposedPorts: map[nat.Port]struct{}{
			port: {},
		},
	}, &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{
			port: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: req.GetPort(),
				},
			},
		},
	}, nil, nil, "")
	s.log.Info("containerId: ", resp.ID)
	if err != nil {
		return nil, err
	}
	if err := s.cli.ContainerStart(ctx2, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}
	return s.GetVm(ctx, &pb.GetVmRequest{
		Id: resp.ID,
	})
}
func (s *VmService) DeleteVm(ctx context.Context, req *pb.DeleteVmRequest) (*pb.DeleteVmReply, error) {
	err := s.cli.ContainerRemove(context.Background(), req.Id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	return &pb.DeleteVmReply{}, err
}
func (s *VmService) GetVm(ctx context.Context, req *pb.GetVmRequest) (*pb.GetVmReply, error) {
	containers, err := s.cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.Arg("id", req.Id)),
	})
	if len(containers) > 0 {
		return toListVmReply(containers)[0], err
	}
	return &pb.GetVmReply{}, err
}
func (s *VmService) ListVm(ctx context.Context, req *pb.ListVmRequest) (*pb.ListVmReply, error) {
	containers, err := s.cli.ContainerList(ctx, types.ContainerListOptions{})
	return &pb.ListVmReply{
		Result: toListVmReply(containers),
	}, err
}

func toListVmReply(containers []types.Container) []*pb.GetVmReply {
	return lo.Map(containers, func(container types.Container, _ int) *pb.GetVmReply {
		return &pb.GetVmReply{
			Id:    container.ID,
			Image: container.Image,
			Ports: lo.Map(container.Ports, func(port types.Port, _ int) *pb.PortBinding {
				return &pb.PortBinding{
					Ip:          port.IP,
					PrivatePort: uint32(port.PrivatePort),
					PublicPort:  uint32(port.PublicPort),
					Type:        port.Type,
				}
			}),
		}
	})
}
