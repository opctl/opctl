package core

import (
	"github.com/opctl/opctl/sdk/go/model"
)

func (this _core) KillOp(
	req model.KillOpReq,
) {
	this.callKiller.Kill(
		req.OpID,
		req.OpID,
	)
}
