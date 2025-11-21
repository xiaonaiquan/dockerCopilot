package utiles

import (
	"context"
	"dockerCopilot/internal/svc"

	"github.com/docker/docker/api/types"
)

func GetContainerInspect(ctx *svc.ServiceContext, id string) (types.ContainerJSON, error) {
	inspectedContainer, err := ctx.DockerClient.ContainerInspect(context.TODO(), id)
	if err != nil {
		return types.ContainerJSON{}, err
	}
	return inspectedContainer, nil
}
