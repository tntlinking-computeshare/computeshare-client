package biz

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	queueTaskV1 "github.com/mohaijiang/computeshare-server/api/queue/v1"
	"io"
	"os"
)

func NewStorageProvider() *StorageProvider {
	return &StorageProvider{
		status: false,
	}
}

type StorageProvider struct {
	status bool
}

func (sp *StorageProvider) Status() bool {
	return sp.status
}

func (sp *StorageProvider) Start(param *queueTaskV1.StorageSetupTaskParamVO) error {
	sp.status = true
	fmt.Println(param.PublicPort)

	return nil
}

func (sp *StorageProvider) Stop() error {
	sp.status = false
	return nil
}

func runSeaweedFSContainer() error {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	imageName := "chrislusf/seaweedfs:3.58"
	containerConfig := &container.Config{
		Image: imageName,
		Cmd: []string{
			"volume",
			"-mserver=computeshare.newtouch.com:9333",
			"-ip.bind=0.0.0.0",
			"-port=41016",
			"-port.public=41016",
			"-ip=computeshare.newtouch.com",
			"-port.grpc=41017",
			"-publicUrl=computeshare.newtouch.com:41016",
		},
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		ExposedPorts: map[nat.Port]struct{}{
			"41016/tcp": struct{}{},
			"41017/tcp": struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"41016/tcp": []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: "41016"},
			},
			"41017/tcp": []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: "41017"},
			},
		},
		Binds: []string{
			"$HOME/.seaweed:/data",
		},
	}

	resp, err := cli.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		"seaweedfs-volume",
	)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return err
	}

	// Print container logs
	io.Copy(os.Stdout, out)

	return nil
}
