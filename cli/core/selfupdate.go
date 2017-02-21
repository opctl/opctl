package core

import (
	"fmt"
	"github.com/opspec-io/opctl/util/cliexiter"
)

func (this _core) SelfUpdate(
	releaseChannel string,
) {

	if releaseChannel != "alpha" && releaseChannel != "beta" && releaseChannel != "stable" {
		this.cliExiter.Exit(
			cliexiter.ExitReq{
				Message: fmt.Sprintf(
					"%v is not an available release channel. "+
						"Available release channels are 'alpha', 'beta', and 'stable'. \n", releaseChannel),
				Code: 1,
			},
		)
		return // support fake exiter
	}

	update, err := this.updater.GetUpdateIfExists(releaseChannel)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	} else if nil == update {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: "No update available, already at the latest version!", Code: 0})
		return // support fake exiter
	}

	err = this.updater.ApplyUpdate(update)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
	} else {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: fmt.Sprintf("Updated to new version: %s!\n", update.Version), Code: 0})
	}
}
