package docker

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func (ctp _containerRuntime) DeleteContainerIfExists(
	ctx context.Context,
	containerID string,
) error {
	err := ctp.dockerClient.ContainerRemove(
		ctx,
		getContainerName(containerID),
		types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		},
	)
	if nil != err {
		return fmt.Errorf("unable to delete container. Response from docker was: %v", err.Error())
	}

	return nil
}
