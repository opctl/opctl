package docker

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"

	dockerClientPkg "github.com/docker/docker/client"
)

func ensureNetworkDetached(
	ctx context.Context,
	dockerClient dockerClientPkg.CommonAPIClient,
) error {
	err := dockerClient.NetworkRemove(ctx, networkName)
	if err != nil && !dockerClientPkg.IsErrNotFound(err) {
		return err
	}

	if runtime.GOOS == "darwin" {
		tunIndex, err := getCurrentTunIndex(ctx)
		if err != nil {
			return err
		}

		tunName := fmt.Sprintf("tun%d", tunIndex)

		cmd := exec.Command("ip", "link", "delete", tunName)

		outputBytes, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Errorf("Failed to delete %s: %w, %s", tunName, err, string(outputBytes))
		}
	}

	return nil
}
