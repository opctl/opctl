package docker

import (
	"fmt"

	"github.com/docker/docker/api/types"
	dockerClientPkg "github.com/docker/docker/client"
	"golang.org/x/net/context"
	"golang.org/x/sync/singleflight"
)

// singleFlightGroup is used to ensure creates don't race across calls
var createSingleFlightGroup singleflight.Group

//counterfeiter:generate -o internal/fakes/ensureNetworkExistser.go . ensureNetworkExistser
type ensureNetworkExistser interface {
	EnsureNetworkExists(
		ctx context.Context,
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
	ctx context.Context,
	networkID string,
) error {

	_, networkInspectErr := ene.dockerClient.NetworkInspect(
		ctx,
		networkID,
		types.NetworkInspectOptions{},
	)
	if networkInspectErr == nil {
		// if network exists, we're done
		return nil
	}

	if !dockerClientPkg.IsErrNotFound(networkInspectErr) {
		return fmt.Errorf("unable to inspect network: %w", networkInspectErr)
	}

	// attempt to resolve within singleFlight.Group to ensure concurrent creates don't race
	_, err, _ := createSingleFlightGroup.Do(
		networkID,
		func() (interface{}, error) {
			return ene.dockerClient.NetworkCreate(
				ctx,
				networkID,
				types.NetworkCreate{
					CheckDuplicate: true,
					Attachable:     true,
				},
			)
		},
	)
	if err != nil {
		return fmt.Errorf("unable to create network: %w", err)
	}

	return nil
}
