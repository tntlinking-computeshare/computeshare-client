package service

import (
	"bytes"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/daemon/logger/jsonfilelog"
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
	"path/filepath"
	"time"

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
	s.log.Info("client开始处理.py脚本，cid: ", req.Cid)
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
	s.log.Info("通过cid获取ipfs资源完成")
	//判断是不是服务器自己部署（/root/client_share_data）
	sharePath := "/root/client_share_data"
	_, err = os.Stat(sharePath)
	currentDir := ""
	if err == nil {
		currentDir = sharePath
	} else if os.IsNotExist(err) {
		currentDir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	} else {
		s.log.Error("判断文件存在不存在失败")
		return nil, err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	// 定义要创建的文件名
	fileName := req.GetCid() + ".py"
	filePath := filepath.Join(currentDir, fileName)
	// 使用 os.Create 创建文件
	create, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	_, err = create.Write(data)
	if err != nil {
		return nil, err
	}
	defer os.Remove(filePath)
	s.log.Info("写入文件成功")
	imageName := "python:3"
	out, err := s.dockerCli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		s.log.Info("拉取镜像失败，err is:", err)
		return nil, err
	}
	s.log.Info("拉取镜像成功，result is:", out)
	//docker执行.py
	var mapping []string
	mapping = append(mapping, filePath+":/tmp/"+fileName)
	var cmd []string
	cmd = append(cmd, "python")
	cmd = append(cmd, "/tmp/"+fileName)
	resp, err := s.dockerCli.ContainerCreate(ctx, &container.Config{
		Image:        imageName,
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{

		Binds:        mapping,
		PortBindings: map[nat.Port][]nat.PortBinding{},
		AutoRemove:   false,
		LogConfig: container.LogConfig{
			Type: jsonfilelog.Name,
		},
	}, nil, nil, "")
	if err != nil {
		return nil, err
	}
	s.log.Info("创建container完成 containerId: ", resp.ID)
	defer s.dockerCli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err := s.dockerCli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		s.log.Error("启动容器失败", err)
		return nil, err
	}
	s.log.Info("container启动成功 containerId: ", resp.ID)
	for {
		inspect, err := s.dockerCli.ContainerInspect(ctx, resp.ID)
		if err != nil {
			s.log.Error("获取容器Inspect失败", err)
			return nil, err
		}
		s.log.Info("container 当前状态是", inspect.State.Status)
		if inspect.State.Status == "exited" {
			logs, err := s.dockerCli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Details: true})
			if err != nil {
				s.log.Error("获取容器Log失败", err)
				return nil, err
			}
			defer logs.Close()
			actualStdout := new(bytes.Buffer)
			actualStderr := io.Discard
			_, err = stdcopy.StdCopy(actualStdout, actualStderr, logs)
			s.log.Info("容器执行的日志是-->", actualStdout.String())
			return &pb.RunPythonPackageClientReply{ExecuteResult: actualStdout.String()}, nil
		} else {
			time.Sleep(time.Millisecond * 500)
		}
	}
}
func (s *ComputePowerService) CancelExecPythonPackage(ctx context.Context, req *pb.CancelExecPythonPackageClientRequest) (*pb.CancelExecPythonPackageClientReply, error) {
	return &pb.CancelExecPythonPackageClientReply{}, nil
}
