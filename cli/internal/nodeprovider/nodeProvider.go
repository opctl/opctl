package nodeprovider

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"github.com/opctl/opctl/cli/internal/model"
)

//counterfeiter:generate -o fakes/nodeProvider.go . NodeProvider
type NodeProvider interface {
	ListNodes() (nodes []model.NodeHandle, err error)
	CreateNodeIfNotExists() (nodeHandle model.NodeHandle, err error)
	KillNodeIfExists(nodeId string) (err error)
}
