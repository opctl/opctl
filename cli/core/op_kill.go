package core

import (
	"context"
	"github.com/opctl/opctl/sdk/go/model"
	"github.com/opctl/opctl/util/cliexiter"
)

func (this _core) OpKill(
	ctx context.Context,
	opId string,
) {

	err := this.apiClient.KillOp(
		ctx,
		model.KillOpReq{
			OpID: opId,
		},
	)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

}
