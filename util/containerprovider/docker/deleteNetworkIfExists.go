package docker

import (
	"fmt"
	"golang.org/x/net/context"
)

func (this _containerProvider) DeleteNetworkIfExists(
	networkId string,
) (err error) {
	err = this.dockerClient.NetworkRemove(
		context.Background(),
		networkId,
	)
	if nil != err {
		err = fmt.Errorf("Unable to delete network. Response from docker was:\n %v", err.Error())
	}
	return
}
