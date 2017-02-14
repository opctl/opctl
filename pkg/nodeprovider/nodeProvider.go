package nodeprovider

import "github.com/opspec-io/opctl/pkg/node"

type NodeProvider interface {
	ListNodes() (nodes []*node.InfoView, err error)
	CreateNode() (nodeId string, err error)
}

type nodeProvider struct{}
