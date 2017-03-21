package core

import "github.com/opctl/opctl/util/cliexiter"

func (this _core) NodeKill() {
	err := this.nodeProvider.KillNodeIfExists("")
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}
