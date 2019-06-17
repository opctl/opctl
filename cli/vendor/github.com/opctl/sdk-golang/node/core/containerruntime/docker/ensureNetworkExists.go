package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func (ctp _containerRuntime) EnsureNetworkExists(
	networkId string,
) (err error) {

	_, networkInspectErr := ctp.dockerClient.NetworkInspect(
		context.Background(),
		networkId,
		types.NetworkInspectOptions{},
	)
	if nil == networkInspectErr {
		// if network exists, we're done
		return
	}

	if !client.IsErrNotFound(networkInspectErr) {
		err = fmt.Errorf("unable to inspect network. Response from docker was:\n %v", networkInspectErr.Error())
		return
	}

	_, err = ctp.dockerClient.NetworkCreate(
		context.Background(),
		networkId,
		types.NetworkCreate{
			CheckDuplicate: true,
			Attachable:     true,
		},
	)
	if nil != err {
		err = fmt.Errorf("unable to create network. Response from docker was:\n %v", err.Error())
	}
	return
}
