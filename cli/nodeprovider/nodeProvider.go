package nodeprovider

import (
	"github.com/opctl/opctl/cli/types"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ NodeProvider

type NodeProvider interface {
	ListNodes() (nodes []*types.NodeInfoView, err error)
	CreateNode() (nodeInfo *types.NodeInfoView, err error)
	KillNodeIfExists(nodeId string) (err error)
}
