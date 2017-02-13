package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func (this _containerProvider) DeleteContainerIfExists(
	containerId string,
) {
	err := this.dockerClient.ContainerRemove(
		context.Background(),
		containerId,
		types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		},
	)
	if nil != err {
		fmt.Printf("Unable to delete container. Response from docker was:\n %v", err.Error())
	}
}
