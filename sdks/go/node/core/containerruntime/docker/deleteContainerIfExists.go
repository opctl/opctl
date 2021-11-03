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
	err := ctp.dockerClient.ContainerStop(
		ctx,
		getContainerName(containerID),
		&stopTimeout,
	)
	if err != nil {
		return fmt.Errorf("unable to stop container: %w", err)
	}

	// now delete the container post-termination
	err = ctp.dockerClient.ContainerRemove(
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
