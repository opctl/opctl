package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
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
		return errors.Wrap(err, "unable to delete container")
	}

	return nil
}
