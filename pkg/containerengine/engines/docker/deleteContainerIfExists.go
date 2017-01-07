package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func (this _containerEngine) DeleteContainerIfExists(
	containerId string,
) {
	err := this.dockerEngine.ContainerRemove(
		context.Background(),
		containerId,
		types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		},
	)
	if nil != err {
		fmt.Printf("Unable to delete container. Response from docker engine was:\n %v", err.Error())
	}
}
