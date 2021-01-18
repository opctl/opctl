package nodeprovider

import (
	"github.com/opctl/opctl/sdks/go/node"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o fakes/nodeProvider.go . NodeProvider
type NodeProvider interface {
	ListNodes() ([]node.Node, error)
	CreateNodeIfNotExists() (node.Node, error)
	KillNodeIfExists(nodeID string) error
}
