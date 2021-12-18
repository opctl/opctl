package docker

import (
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func (ctp _containerRuntime) DeleteContainerIfExists(
	ctx context.Context,
	containerID string,
) error {
	// try to stop the container gracefully prior to deletion
	stopTimeout := 3 * time.Second
	// ignore error; we want to remove regardless
	ctp.dockerClient.ContainerStop(
		ctx,
		getContainerName(containerID),
		&stopTimeout,
	)

	// now delete the container post-termination
	err := ctp.dockerClient.ContainerRemove(
		ctx,
		getContainerName(containerID),
		types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		},
	)
	if err != nil {
		return fmt.Errorf("unable to delete container: %w", err)
	}

	return nil
}
