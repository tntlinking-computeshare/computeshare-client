package biz

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-kratos/kratos/v2/log"
	queueTaskV1 "github.com/mohaijiang/computeshare-server/api/queue/v1"
	"os"
	"strconv"
)

const seaweedContainerName = "seaweedfs-volume"

func NewStorageProvider(logger log.Logger) *StorageProvider {
	return &StorageProvider{
		log: log.NewHelper(logger),
	}
}

type StorageProvider struct {
	log *log.Helper
}

func (sp *StorageProvider) Status() bool {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		sp.log.Error(err)
		return false
	}

	inspect, err := cli.ContainerInspect(ctx, seaweedContainerName)
	if err != nil {
		sp.log.Error(err)
		return false
	}
	return inspect.State.Running
}

func (sp *StorageProvider) Start(param *queueTaskV1.StorageSetupTaskParamVO) error {

	return sp.runSeaweedFSContainer(param)
}

func (sp *StorageProvider) Stop() error {
	return sp.stopSeaweedFSContainer()
}

// runSeaweedFSContainer
func (sp *StorageProvider) runSeaweedFSContainer(param *queueTaskV1.StorageSetupTaskParamVO) error {
	sp.log.Info("runSeaweedFSContainer")
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		sp.log.Error(err)
		return err
	}

	inspect, err := cli.ContainerInspect(ctx, seaweedContainerName)
	if err == nil && inspect.State.Running {
		sp.log.Error(fmt.Sprintf("容器%s已经在运行,id: %s\n", inspect.Name, inspect.ID))
		return errors.New(fmt.Sprintf("容器%s已经在运行,id: %s\n", inspect.Name, inspect.ID))
	}

	if err == nil && !inspect.State.Running {
		err = cli.ContainerRemove(ctx, seaweedContainerName, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			sp.log.Error(err)
			return err
		}
	}

	imageName := "chrislusf/seaweedfs:3.58"
	containerConfig := &container.Config{
		Image: imageName,
		Cmd: []string{
			"volume",
			fmt.Sprintf("-mserver=%s", param.MasterServer),
			"-ip.bind=0.0.0.0",
			fmt.Sprintf("-port=%d", param.PublicPort),
			fmt.Sprintf("-port.public=%d", param.PublicPort),
			fmt.Sprintf("-ip=%s", param.PublicIp),
			fmt.Sprintf("-port.grpc=%d", param.GrpcPort),
			fmt.Sprintf("-publicUrl=%s:%d", param.PublicIp, param.PublicPort),
		},
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		ExposedPorts: map[nat.Port]struct{}{
			nat.Port(fmt.Sprintf("%d/tcp", param.PublicPort)): struct{}{},
			nat.Port(fmt.Sprintf("%d/tcp", param.GrpcPort)):   struct{}{},
		},
	}

	home, _ := os.UserHomeDir()

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(fmt.Sprintf("%d/tcp", param.PublicPort)): []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: strconv.Itoa(int(param.PublicPort))},
			},
			nat.Port(fmt.Sprintf("%d/tcp", param.GrpcPort)): []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: strconv.Itoa(int(param.GrpcPort))},
			},
		},
		Binds: []string{
			fmt.Sprintf("%s/.seaweed:/data", home),
		},
	}

	resp, err := cli.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		seaweedContainerName,
	)
	if err != nil {
		sp.log.Error(err)
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		sp.log.Error(err)
		return err
	}

	return nil
}

func (sp *StorageProvider) stopSeaweedFSContainer() error {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		sp.log.Error(err)
		return nil
	}

	inspect, err := cli.ContainerInspect(ctx, seaweedContainerName)
	if err != nil {
		fmt.Println("容器已经停止")
		sp.log.Error(err)
		return err
	}

	if inspect.State.Running {
		err = cli.ContainerStop(ctx, seaweedContainerName, container.StopOptions{})
		if err != nil {
			sp.log.Error(err)
			return err
		}
	}

	err = cli.ContainerRemove(ctx, seaweedContainerName, types.ContainerRemoveOptions{
		Force: true,
	})

	return err
}
