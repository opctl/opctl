package docker

import (
	"fmt"
	"github.com/docker/docker/api/types/network"
	"golang.org/x/net/context"
)

func (this _containerProvider) NetworkContainer(
	networkId string,
	containerId string,
	containerAlias string,
) (err error) {
	err = this.dockerClient.NetworkConnect(
		context.Background(),
		networkId,
		containerId,
		&network.EndpointSettings{
			Aliases: []string{containerAlias},
		},
	)
	if nil != err {
		err = fmt.Errorf("Unable to network container. Response from docker was:\n %v", err.Error())
	}
	return
}
