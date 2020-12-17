package core

import (
	"github.com/opctl/opctl/cli/internal/core/node"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

// Noder exposes the "node" sub command
type Noder interface {
	Node() node.Node
}

// newNoder returns an initialized "node" sub command
func newNoder(nodeProvider nodeprovider.NodeProvider) Noder {
	return _noder{
		node: node.New(nodeProvider),
	}
}

type _noder struct {
	node node.Node
}

func (ivkr _noder) Node() node.Node {
	return ivkr.node
}
