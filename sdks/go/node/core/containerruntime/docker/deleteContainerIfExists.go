package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func (ctp _containerRuntime) DeleteContainerIfExists(
	containerName string,
) (err error) {
	err = ctp.dockerClient.ContainerRemove(
		context.Background(),
		containerName,
		types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		},
	)
	if nil != err {
		err = fmt.Errorf("unable to delete container. Response from docker was: %v", err.Error())
	}
	return
}
