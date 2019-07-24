package core

import (
	"github.com/opctl/opctl/sdks/go/types"
)

func (this _core) KillOp(
	req types.KillOpReq,
) {
	this.callKiller.Kill(
		req.OpID,
		req.OpID,
	)
}
