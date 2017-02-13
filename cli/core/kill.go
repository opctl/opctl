package core

import (
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

func (this _core) KillOp(
	opGraphId string,
) {

	err := this.engineClient.KillOp(
		model.KillOpReq{
			OpGraphId: opGraphId,
		},
	)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

}
