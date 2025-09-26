package vm

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	tc "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func deleteContainerByName(ctx context.Context, cli *client.Client, containerName string) error {
	// 创建过滤器查找容器
	filter := filters.NewArgs()
	filter.Add("name", containerName)

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filter,
	})
	if err != nil {
		return fmt.Errorf("获取容器列表失败: %v", err)
	}

	if len(containers) == 0 {
		return fmt.Errorf("容器 %s 不存在", containerName)
	}

	// 删除找到的容器
	for _, container := range containers {
		// 检查容器名是否完全匹配（docker会返回包含该名称的所有容器）
		for _, name := range container.Names {
			if name == "/"+containerName || name == containerName {
				fmt.Printf("正在删除容器: %s (ID: %s)\n", containerName, container.ID[:12])

				// 先停止容器（如果正在运行）
				if container.State == "running" {
					if err := cli.ContainerStop(ctx, container.ID, tc.StopOptions{}); err != nil {
						return fmt.Errorf("停止容器失败: %v", err)
					}
				}

				// 删除容器
				if err := cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{
					Force: true, // 强制删除
				}); err != nil {
					return fmt.Errorf("删除容器失败: %v", err)
				}

				fmt.Printf("成功删除容器: %s\n", containerName)
				return nil
			}
		}
	}

	return fmt.Errorf("未找到完全匹配的容器: %s", containerName)
}
