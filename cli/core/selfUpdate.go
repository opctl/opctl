package core

import (
	"fmt"
)

func (this _core) SelfUpdate(
	releaseChannel string,
) {

	if releaseChannel != "beta" && releaseChannel != "stable" {
		this.exiter.Exit(
			ExitReq{
				Message: fmt.Sprintf(
					"%v is not an available release channel. "+
						"Available release channels are 'beta' 'stable'. \n", releaseChannel),
				Code: 1,
			},
		)
		return // support fake exiter
	}

	update, err := this.updater.TryGetUpdate(releaseChannel)
	if nil != err {
		this.exiter.Exit(ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	} else if nil == update {
		this.exiter.Exit(ExitReq{Message: "No update available, already at the latest version!", Code: 0})
		return // support fake exiter
	}

	err = this.updater.ApplyUpdate(update)
	if nil != err {
		this.exiter.Exit(ExitReq{Message: err.Error(), Code: 1})
	} else {
		this.exiter.Exit(ExitReq{Message: fmt.Sprintf("Updated to new version: %s!\n", update.Version), Code: 0})
	}
}
