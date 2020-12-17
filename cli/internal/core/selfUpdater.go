package core

import (
	"fmt"

	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/cli/internal/updater"
)

// SelfUpdater exposes the "self-update" command
type SelfUpdater interface {
	SelfUpdate(
		releaseChannel string,
	) (string, error)
}

// newSelfUpdater returns an initialized "self-update" command
func newSelfUpdater(nodeProvider nodeprovider.NodeProvider) SelfUpdater {
	return _selfUpdateInvoker{
		nodeProvider: nodeProvider,
		updater:      updater.New(),
	}
}

type _selfUpdateInvoker struct {
	nodeProvider nodeprovider.NodeProvider
	updater      updater.Updater
}

func (ivkr _selfUpdateInvoker) SelfUpdate(
	releaseChannel string,
) (string, error) {
	if releaseChannel != "alpha" && releaseChannel != "beta" && releaseChannel != "stable" {
		return "", fmt.Errorf(
			"%v is not an available release channel. "+
				"Available release channels are 'alpha', 'beta', and 'stable'.", releaseChannel)
	}

	update, err := ivkr.updater.GetUpdateIfExists(releaseChannel)
	if nil != err {
		return "", err
	} else if nil == update {
		return "No update available, already at the latest version!", err
	}

	err = ivkr.updater.ApplyUpdate(update)
	if nil != err {
		return "", err
	}

	// kill local node to ensure outdated version not left running
	err = ivkr.nodeProvider.KillNodeIfExists("")
	if nil != err {
		return "", fmt.Errorf("Unable to kill running node; run `node kill` to complete the update. Error was: %v", err)
	}

	// @TODO start node maintaining previous user
	return fmt.Sprintf("Updated to new version: %s!", update.Version), nil
}
