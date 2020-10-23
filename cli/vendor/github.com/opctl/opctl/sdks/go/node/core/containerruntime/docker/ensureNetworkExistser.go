package docker

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	dockerClientPkg "github.com/docker/docker/client"
	"golang.org/x/net/context"
	"golang.org/x/sync/singleflight"
)

// singleFlightGroup is used to ensure creates don't race across calls
var createSingleFlightGroup singleflight.Group

//counterfeiter:generate -o internal/fakes/ensureNetworkExistser.go . ensureNetworkExistser
type ensureNetworkExistser interface {
	EnsureNetworkExists(
		networkID string,
	) error
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
) error {

	_, networkInspectErr := ene.dockerClient.NetworkInspect(
		context.Background(),
		networkID,
		types.NetworkInspectOptions{},
	)
	if nil == networkInspectErr {
		// if network exists, we're done
		return nil
	}

	if !client.IsErrNotFound(networkInspectErr) {
		return fmt.Errorf("unable to inspect network. Response from docker was: %v", networkInspectErr.Error())
	}

	// attempt to resolve within singleFlight.Group to ensure concurrent creates don't race
	_, err, _ := createSingleFlightGroup.Do(
		networkID,
		func() (interface{}, error) {
			return ene.dockerClient.NetworkCreate(
				context.Background(),
				networkID,
				types.NetworkCreate{
					CheckDuplicate: true,
					Attachable:     true,
				},
			)
		},
	)
	if nil != err {
		return fmt.Errorf("unable to create network. Response from docker was: %v", err.Error())
	}

	return nil
}
