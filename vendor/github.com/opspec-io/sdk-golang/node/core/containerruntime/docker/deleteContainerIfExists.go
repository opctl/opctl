package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func (ctp _containerRuntime) DeleteContainerIfExists(
	containerID string,
) (err error) {
	err = ctp.dockerClient.ContainerRemove(
		context.Background(),
		containerID,
		types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		},
	)
	if nil != err {
		err = fmt.Errorf("unable to delete container. Response from docker was:\n %v", err.Error())
	}
	return
}
