package main

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

func selfUpdate(
	nodeProvider nodeprovider.NodeProvider,
) (string, error) {
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
	err = nodeProvider.KillNodeIfExists("")
	if err != nil {
		err = fmt.Errorf("unable to kill running node; run `node kill` to complete the update: %w", err)
	}
	return fmt.Sprintf("Updated to new version: %s!", latest.Version), err
}
