package core

import (
	"context"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/sdk-golang/model"
)

func (this _core) OpKill(
	ctx context.Context,
	opId string,
) {

	err := this.opspecNodeAPIClient.KillOp(
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
