package docker

import (
	"fmt"
	"strings"

	"context"

	"github.com/docker/docker/api/types"
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/opctl/opctl/sdks/go/model"
)

func ensureNetworkExists(
	ctx context.Context,
	dockerClient dockerClientPkg.CommonAPIClient,
	imagePullCreds *model.Creds,
	networkName string,
) error {
	// always attempt to create to avoid races
	_, networkCreateErr := dockerClient.NetworkCreate(
		ctx,
		networkName,
		types.NetworkCreate{
			CheckDuplicate: true,
			Attachable:     true,
		},
	)
	// return errors not related to already existing...
	if networkCreateErr != nil && !strings.Contains(networkCreateErr.Error(), "exists") {
		return fmt.Errorf("unable to create network: %w", networkCreateErr)
	}

	return ensureNetworkAttached(
		ctx,
		dockerClient,
		imagePullCreds,
	)
}
