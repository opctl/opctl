package core

import (
	"fmt"

	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/cli/internal/updater"
)

// SelfUpdater exposes the "self-update" command
type SelfUpdater interface {
	SelfUpdate(
		releaseChannel string,
	)
}

// newSelfUpdater returns an initialized "self-update" command
func newSelfUpdater(
	cliExiter cliexiter.CliExiter,
	nodeProvider nodeprovider.NodeProvider,
) SelfUpdater {
	return _selfUpdateInvoker{
		cliExiter:    cliExiter,
		nodeProvider: nodeProvider,
		updater:      updater.New(),
	}
}

type _selfUpdateInvoker struct {
	cliExiter    cliexiter.CliExiter
	nodeProvider nodeprovider.NodeProvider
	updater      updater.Updater
}

func (ivkr _selfUpdateInvoker) SelfUpdate(
	releaseChannel string,
) {

	if releaseChannel != "alpha" && releaseChannel != "beta" && releaseChannel != "stable" {
		ivkr.cliExiter.Exit(
			cliexiter.ExitReq{
				Message: fmt.Sprintf(
					"%v is not an available release channel. "+
						"Available release channels are 'alpha', 'beta', and 'stable'. \n", releaseChannel),
				Code: 1,
			},
		)
		return // support fake exiter
	}

	update, err := ivkr.updater.GetUpdateIfExists(releaseChannel)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{
			Message: err.Error(),
			Code:    1,
		})
		return // support fake exiter
	} else if nil == update {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{
			Message: "No update available, already at the latest version!",
			Code:    0,
		})
		return // support fake exiter
	}

	err = ivkr.updater.ApplyUpdate(update)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	// kill local node to ensure outdated version not left running
	err = ivkr.nodeProvider.KillNodeIfExists("")
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{
			Message: fmt.Sprintf("Unable to kill running node; run `node kill` to complete the update. Error was: %v", err.Error()),
			Code:    1,
		})
		return // support fake exiter
	}

	// @TODO start node maintaining previous user

	ivkr.cliExiter.Exit(cliexiter.ExitReq{
		Message: fmt.Sprintf("Updated to new version: %s!\n", update.Version),
		Code:    0,
	})

}
