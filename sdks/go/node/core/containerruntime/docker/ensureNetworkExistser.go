package docker

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	dockerClientPkg "github.com/docker/docker/client"
	"golang.org/x/net/context"
)

//counterfeiter:generate -o internal/fakes/ensureNetworkExistser.go . ensureNetworkExistser
type ensureNetworkExistser interface {
	EnsureNetworkExists(
		networkID string,
	) (err error)
}

func newEnsureNetworkExistser(dockerClient dockerClientPkg.CommonAPIClient) ensureNetworkExistser {
	return _ensureNetworkExistser{
		dockerClient: dockerClient,
	}
}

type _ensureNetworkExistser struct {
	dockerClient dockerClientPkg.CommonAPIClient
}

func (ene _ensureNetworkExistser) EnsureNetworkExists(
	networkID string,
) (err error) {

	_, networkInspectErr := ene.dockerClient.NetworkInspect(
		context.Background(),
		networkID,
		types.NetworkInspectOptions{},
	)
	if nil == networkInspectErr {
		// if network exists, we're done
		return
	}

	if !client.IsErrNotFound(networkInspectErr) {
		err = fmt.Errorf("unable to inspect network. Response from docker was: %v", networkInspectErr.Error())
		return
	}

	_, err = ene.dockerClient.NetworkCreate(
		context.Background(),
		networkID,
		types.NetworkCreate{
			CheckDuplicate: true,
			Attachable:     true,
		},
	)
	if nil != err {
		err = fmt.Errorf("unable to create network. Response from docker was: %v", err.Error())
	}
	return
}
