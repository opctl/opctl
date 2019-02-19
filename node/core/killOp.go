package core

import (
	"github.com/opctl/sdk-golang/model"
)

func (this _core) KillOp(
	req model.KillOpReq,
) {
	this.opKiller.Kill(req)
}
