package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"
	pb "github.com/mohaijiang/computeshare-client/api/compute/v1"
	"github.com/samber/lo"
	"io"
	"net/http"
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
	out, err := s.cli.ImagePull(ctx, req.Image, types.ImagePullOptions{})
	s.log.Info(out)
	if err != nil {
		return nil, err
	}

	port, err := nat.NewPort("tcp", req.GetPort())
	resp, err := s.cli.ContainerCreate(ctx, &container.Config{
		Image: req.Image,
		ExposedPorts: map[nat.Port]struct{}{
			port: {},
		},
		Cmd: req.Command,
	}, &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{
			//port: []nat.PortBinding{
			//	{
			//		HostIP:   "0.0.0.0",
			//		HostPort: req.GetPort(),
			//	},
			//},
		},
	}, nil, nil, "")
	s.log.Info("containerId: ", resp.ID)
	if err != nil {
		return nil, err
	}
	if err := s.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
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

	result := &pb.GetVmReply{
		Id: req.GetId(),
	}

	containerJson, err := s.cli.ContainerInspect(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	result.Image = containerJson.Image

	stats, err := s.cli.ContainerStats(ctx, req.GetId(), true)
	if err == nil {
		defer stats.Body.Close()

		// 解析容器统计信息
		var statsData types.Stats
		if err := json.NewDecoder(stats.Body).Decode(&statsData); err != nil {
			log.Fatal(err)
		}

		// 从 statsData 中提取 CPU 和内存使用情况
		cpuUsage := statsData.CPUStats.CPUUsage.TotalUsage
		memoryUsage := statsData.MemoryStats.Usage

		// 打印 CPU 和内存使用情况
		fmt.Printf("CPU Usage: %d\n", cpuUsage)
		fmt.Printf("Memory Usage: %d\n", memoryUsage)
		result.CpuUsage = cpuUsage
		result.MemoryUsage = memoryUsage
	}

	return result, err
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

func (s *VmService) StartVm(ctx context.Context, req *pb.GetVmRequest) (*pb.GetVmReply, error) {
	err := s.cli.ContainerStart(ctx, req.GetId(), types.ContainerStartOptions{})
	return &pb.GetVmReply{}, err
}
func (s *VmService) StopVm(ctx context.Context, req *pb.GetVmRequest) (*pb.GetVmReply, error) {
	timeout := 2
	err := s.cli.ContainerStop(ctx, req.GetId(), container.StopOptions{
		Timeout: &timeout,
	})

	return &pb.GetVmReply{}, err
}

type VmWebsocketHandler struct {
	cli *client.Client
	ctx context.Context
}

func NewVmWebsocketHandler(cli *client.Client) *VmWebsocketHandler {
	return &VmWebsocketHandler{
		cli: cli,
		ctx: context.Background(),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (handler *VmWebsocketHandler) Terminal(w http.ResponseWriter, r *http.Request) {
	// websocket握手
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()

	r.ParseForm()
	// 获取容器ID或name
	container := r.Form.Get("container")
	// 执行exec，获取到容器终端的连接
	hr, err := handler.exec(container, r.Form.Get("workdir"))
	if err != nil {
		log.Error(err)
		return
	}
	// 关闭I/O流
	defer hr.Close()
	// 退出进程
	defer func() {
		hr.Conn.Write([]byte("exit\r"))
	}()

	go func() {
		handler.wsWriterCopy(hr.Conn, conn)
	}()
	handler.wsReaderCopy(conn, hr.Conn)
}

func (handler *VmWebsocketHandler) exec(container string, workdir string) (hr types.HijackedResponse, err error) {
	// 执行/bin/bash命令
	ir, err := handler.cli.ContainerExecCreate(handler.ctx, container, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   workdir,
		Cmd:          []string{"/bin/bash"},
		Tty:          true,
	})
	if err != nil {
		return
	}

	// 附加到上面创建的/bin/bash进程中
	hr, err = handler.cli.ContainerExecAttach(handler.ctx, ir.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		return
	}
	return
}

func (handler *VmWebsocketHandler) wsWriterCopy(reader io.Reader, writer *websocket.Conn) {
	buf := make([]byte, 8192)
	for {
		nr, err := reader.Read(buf)
		if nr > 0 {
			err := writer.WriteMessage(websocket.BinaryMessage, buf[0:nr])
			if err != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}
}

func (handler *VmWebsocketHandler) wsReaderCopy(reader *websocket.Conn, writer io.Writer) {
	for {
		messageType, p, err := reader.ReadMessage()
		if err != nil {
			return
		}
		if messageType == websocket.TextMessage {
			writer.Write(p)
		}
	}
}
