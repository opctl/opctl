package docker

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/docker/docker/api/types"
	dockerClientPkg "github.com/docker/docker/client"
)

func ensureNetworkDetached(
	ctx context.Context,
	dockerClient dockerClientPkg.CommonAPIClient,
) error {

	if runtime.GOOS == "darwin" {
		networkResource, networkInspectErr := dockerClient.NetworkInspect(
			ctx,
			networkName,
			types.NetworkInspectOptions{},
		)
		if networkInspectErr != nil {
			return fmt.Errorf("unable to inspect network: %w", networkInspectErr)
		}

		for _, config := range networkResource.IPAM.Config {
			if networkResource.Scope == "local" {
				cmd := exec.Command("route", "-q", "-n", "delete", "-inet", config.Subnet)

				outputBytes, err := cmd.CombinedOutput()

				if err != nil {
					fmt.Errorf("Failed to delete route: %w, %s", err, string(outputBytes))
				}
			}
		}
	}

	return nil
}
