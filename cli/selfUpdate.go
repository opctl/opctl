package main

import (
	"fmt"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/cli/internal/updater"
)

func selfUpdate(
	nodeProvider nodeprovider.NodeProvider,
	channel string,
) (string, error) {
	updater := updater.New()
	update, err := updater.GetUpdateIfExists(channel)
	if nil != err {
		return "", err
	} else if nil == update {
		return "No update available, already at the latest version!", nil
	}

	err = updater.ApplyUpdate(update)
	if nil != err {
		return "", err
	}

	// kill local node to ensure outdated version not left running
	// @TODO start node maintaining previous user
	err = nodeProvider.KillNodeIfExists("")
	if nil != err {
		err = fmt.Errorf("Unable to kill running node; run `node kill` to complete the update. Error was: %v", err)
	}
	return fmt.Sprintf("Updated to new version: %s!", update.Version), err
}
