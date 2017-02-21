package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func (this _containerProvider) CreateNetwork(
	networkId string,
) (err error) {
	_, err = this.dockerClient.NetworkCreate(
		context.Background(),
		networkId,
		types.NetworkCreate{
			Attachable: true,
		},
	)
	if nil != err {
		err = fmt.Errorf("Unable to create network. Response from docker was:\n %v", err.Error())
	}
	return
}
