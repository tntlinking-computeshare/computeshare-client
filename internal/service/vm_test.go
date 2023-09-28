package service

import (
	"context"
	"fmt"
	"testing"
)

func TestCreateContainer(t *testing.T) {

	cli, _ := NewDockerCli()

	containerJson, err := cli.ContainerInspect(context.Background(), "aaa")

	fmt.Println(err)
	fmt.Println(containerJson.ContainerJSONBase)
}
