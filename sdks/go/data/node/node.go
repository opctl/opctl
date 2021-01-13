package node

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
)

// New returns a data provider which sources pkgs from a node
func New(
	core node.OpNode,
	pullCreds *model.Creds,
) model.DataProvider {
	return _node{
		core:      core,
		pullCreds: pullCreds,
	}
}

type _node struct {
	core      node.OpNode
	pullCreds *model.Creds
}

func (np _node) TryResolve(
	ctx context.Context,
	dataRef string,
) (model.DataHandle, error) {

	// ensure resolvable by listing contents w/out err
	if _, err := np.core.ListDescendants(
		ctx,
		model.ListDescendantsReq{
			PkgRef:    dataRef,
			PullCreds: np.pullCreds,
		},
	); nil != err {
		return nil, err
	}

	return newHandle(np.core, dataRef, np.pullCreds), nil
}
