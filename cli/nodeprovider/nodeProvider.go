package nodeprovider

import (
	"github.com/opctl/opctl/cli/model"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fake.go --fake-name Fake ./ NodeProvider

type NodeProvider interface {
	ListNodes() (nodes []*model.NodeInfoView, err error)
	CreateNode() (nodeInfo *model.NodeInfoView, err error)
	KillNodeIfExists(nodeId string) (err error)
}
