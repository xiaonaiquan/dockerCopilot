package utiles

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"dockerCopilot/internal/svc"
)

// DeleteContainer 删除容器
func DeleteContainer(ctx *svc.ServiceContext, id string, force bool) error {
	 removeOptions := container.RemoveOptions{RemoveVolumes:false,RemoveLinks:false,Force: force}
	 return ctx.DockerClient.ContainerRemove(context.Background(), id, removeOptions)
}
