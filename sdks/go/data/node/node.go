package node

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
)

// New returns a data provider which sources pkgs from a node
func New(
	node node.Node,
	pullCreds *model.Creds,
) model.DataProvider {
	return _node{
		node:      node,
		pullCreds: pullCreds,
	}
}

type _node struct {
	node      node.Node
	pullCreds *model.Creds
}

func (np _node) Label() string {
	return "opctl node"
}

func (np _node) TryResolve(
	ctx context.Context,
	dataRef string,
) (model.DataHandle, error) {

	// ensure resolvable by listing contents w/out err
	if _, err := np.node.ListDescendants(
		ctx,
		model.ListDescendantsReq{
			PkgRef:    dataRef,
			PullCreds: np.pullCreds,
		},
	); err != nil {
		return nil, err
	}

	return newHandle(np.node, dataRef, np.pullCreds), nil
}
