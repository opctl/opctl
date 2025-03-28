package main

import (
	"context"
	"fmt"

	"github.com/blang/semver"
	"github.com/opctl/opctl/cli/internal/euid0"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

func selfUpdate(
	ctx context.Context,
	nodeConfig local.NodeConfig,
) (string, error) {
	if err := euid0.Ensure(); err != nil {
		return "", err
	}

	v := semver.MustParse(version)
	latest, err := selfupdate.UpdateSelf(v, "opctl/opctl")
	if err != nil {
		return "", err
	}

	if latest.Version.Equals(v) {
		return "No update available, already at the latest version!", nil
	}

	// kill local node to ensure outdated version not left running
	// @TODO start node maintaining previous user
	np, err := local.New(nodeConfig)
	if err != nil {
		return "", err
	}

	err = np.KillNodeIfExists(
		ctx,
	)
	if err != nil {
		err = fmt.Errorf("unable to kill running node; run `node kill` to complete the update: %w", err)
	}
	return fmt.Sprintf("Updated to new version: %s!", latest.Version), err
}
