package op

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
)

// Killer exposes the "op kill" sub command
type Killer interface {
	Kill(
		ctx context.Context,
		opID string,
	) error
}

// newKiller returns an initialized "op kill" sub command
func newKiller(core node.OpNode) Killer {
	return _killer{
		core: core,
	}
}

type _killer struct {
	core node.OpNode
}

func (ivkr _killer) Kill(
	ctx context.Context,
	opID string,
) error {
	return ivkr.core.KillOp(
		ctx,
		model.KillOpReq{
			OpID:       opID,
			RootCallID: opID,
		},
	)
}
