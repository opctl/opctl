package core

import (
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/model"
)

func (this _core) OpKill(
	opId string,
) {

	err := this.consumeNodeApi.KillOp(
		model.KillOpReq{
			OpId: opId,
		},
	)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

}
