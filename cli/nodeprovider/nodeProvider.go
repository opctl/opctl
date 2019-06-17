package nodeprovider

import (
	"github.com/opctl/opctl/cli/model"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ NodeProvider

type NodeProvider interface {
	ListNodes() (nodes []*model.NodeInfoView, err error)
	CreateNode() (nodeInfo *model.NodeInfoView, err error)
	KillNodeIfExists(nodeId string) (err error)
}
