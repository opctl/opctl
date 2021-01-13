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
func newKiller(opNode node.OpNode) Killer {
	return _killer{
		opNode: opNode,
	}
}

type _killer struct {
	opNode node.OpNode
}

func (ivkr _killer) Kill(
	ctx context.Context,
	opID string,
) error {
	return ivkr.opNode.KillOp(
		ctx,
		model.KillOpReq{
			OpID:       opID,
			RootCallID: opID,
		},
	)
}
