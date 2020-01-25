package nodeprovider

import (
	"github.com/opctl/opctl/cli/internal/model"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fake.go --fake-name Fake ./ NodeProvider

type NodeProvider interface {
	ListNodes() (nodes []model.NodeHandle, err error)
	CreateNodeIfNotExists() (nodeHandle model.NodeHandle, err error)
	KillNodeIfExists(nodeId string) (err error)
}
