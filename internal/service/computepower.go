package service

import (
	"bytes"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/go-kratos/kratos/v2/log"
	iface "github.com/ipfs/boxo/coreiface"
	"github.com/ipfs/boxo/coreiface/options"
	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	"io"
	"os"

	pb "github.com/mohaijiang/computeshare-client/api/compute/v1"
)

type ComputePowerService struct {
	pb.UnimplementedComputePowerClientServer
	ipfsNode  *core.IpfsNode
	ipfsApi   iface.CoreAPI
	dockerCli *client.Client
	log       *log.Helper
}

func NewComputePowerService(ipfsNode *core.IpfsNode, client *client.Client, logger log.Logger) (*ComputePowerService, error) {
	api, err := coreapi.NewCoreAPI(ipfsNode, options.Api.FetchBlocks(true))
	if err != nil {
		return nil, err
	}
	return &ComputePowerService{
		ipfsNode:  ipfsNode,
		ipfsApi:   api,
		dockerCli: client,
		log:       log.NewHelper(logger),
	}, nil
}

func (s *ComputePowerService) RunPythonPackage(ctx context.Context, req *pb.RunPythonPackageClientRequest) (*pb.RunPythonPackageClientReply, error) {
	f, err := s.ipfsApi.Unixfs().Get(ctx, path.New(req.Cid))
	var file files.File
	switch f := f.(type) {
	case files.File:
		file = f
	case files.Directory:
		return nil, iface.ErrIsDir
	default:
		return nil, iface.ErrNotSupported
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// 定义要创建的文件名
	fileName := req.GetCid() + ".py"
	// 使用 os.Create 创建文件
	create, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	_, err = create.Write(data)
	if err != nil {
		return nil, err
	}
	// 获取文件的绝对路径
	filePath := currentDir + "/" + fileName

	//docker执行.py
	var mapping []string
	mapping = append(mapping, filePath+":/tmp/"+fileName)
	var cmd []string
	cmd = append(cmd, "python")
	cmd = append(cmd, fileName)
	resp, err := s.dockerCli.ContainerCreate(ctx, &container.Config{
		Image: "python:3",
		Cmd:   cmd,
	}, &container.HostConfig{
		Binds:        mapping,
		PortBindings: map[nat.Port][]nat.PortBinding{},
		AutoRemove:   true,
	}, nil, nil, "")
	s.log.Info("containerId: ", resp.ID)
	if err != nil {
		return nil, err
	}
	if err := s.dockerCli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}
	logs, err := s.dockerCli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return nil, err
	}
	actualStdout := new(bytes.Buffer)
	actualStderr := io.Discard
	_, err = stdcopy.StdCopy(actualStdout, actualStderr, logs)
	return &pb.RunPythonPackageClientReply{ExecuteResult: actualStdout.String()}, nil
}
func (s *ComputePowerService) CancelExecPythonPackage(ctx context.Context, req *pb.CancelExecPythonPackageClientRequest) (*pb.CancelExecPythonPackageClientReply, error) {
	return &pb.CancelExecPythonPackageClientReply{}, nil
}
