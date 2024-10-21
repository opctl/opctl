package docker

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/docker/docker/api/types"
	dockerClientPkg "github.com/docker/docker/client"
)

func ensureNetworkRoutable(
	ctx context.Context,
	dockerClient dockerClientPkg.CommonAPIClient,
) error {

	// Routes to container IPs do not exist on docker for mac. This is inconsistent with linux docker
	// and is so important we implement this bandaid.
	//
	// note: This will likely be brittle because It's relying on undocumented docker for mac internals.
	if runtime.GOOS == "darwin" {
		networkResource, networkInspectErr := dockerClient.NetworkInspect(
			ctx,
			networkName,
			types.NetworkInspectOptions{},
		)
		if networkInspectErr != nil {
			return fmt.Errorf("unable to inspect network: %w", networkInspectErr)
		}

		for _, ipamConfig := range networkResource.IPAM.Config {
			cmd := exec.Command("route", "add", ipamConfig.Subnet, "192.168.64.2")

			// Execute the command
			output, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("unable to ensure network routable %w: %s", err, string(output))
			}
		}
	}

	return nil
}
