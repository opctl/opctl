package docker

import (
	"fmt"
	"strings"

	"context"

	"github.com/docker/docker/api/types/network"
	dockerClientPkg "github.com/docker/docker/client"
)

func ensureNetworkExists(
	ctx context.Context,
	dockerClient dockerClientPkg.CommonAPIClient,
	networkName string,
) error {
	// always attempt to create to avoid races
	_, networkCreateErr := dockerClient.NetworkCreate(
		ctx,
		networkName,
		network.CreateOptions{
			Attachable: true,
		},
	)
	// return errors not related to already existing...
	if networkCreateErr != nil && !strings.Contains(networkCreateErr.Error(), "exists") {
		return fmt.Errorf("unable to create network: %w", networkCreateErr)
	}

	return ensureNetworkAttached(
		ctx,
		dockerClient,
	)
}
