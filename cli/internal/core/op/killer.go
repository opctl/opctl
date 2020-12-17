package op

import (
	"context"

	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/sdks/go/model"
)

// Killer exposes the "op kill" sub command
type Killer interface {
	Kill(
		ctx context.Context,
		opID string,
	) error
}

// newKiller returns an initialized "op kill" sub command
func newKiller(nodeProvider nodeprovider.NodeProvider) Killer {
	return _killer{
		nodeProvider: nodeProvider,
	}
}

type _killer struct {
	nodeProvider nodeprovider.NodeProvider
}

func (ivkr _killer) Kill(
	ctx context.Context,
	opID string,
) error {
	nodeHandle, err := ivkr.nodeProvider.CreateNodeIfNotExists()
	if nil != err {
		return err
	}

	return nodeHandle.APIClient().KillOp(
		ctx,
		model.KillOpReq{
			OpID:       opID,
			RootCallID: opID,
		},
	)
}
