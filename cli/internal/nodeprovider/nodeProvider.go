package nodeprovider

import (
	"context"

	"github.com/opctl/opctl/sdks/go/node"
)

type NodeProvider interface {
	CreateNodeIfNotExists(ctx context.Context) (node.Node, error)
	KillNodeIfExists() error
}
