package containerengine

//go:generate counterfeiter -o engines/fake/containerEngine.go --fake-name ContainerEngine ./ ContainerEngine

import (
	"github.com/opspec-io/opctl/util/eventbus"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

type ContainerEngine interface {
	StartContainer(
		cmd []string,
		env []*model.ContainerInstanceEnvEntry,
		fs []*model.ContainerInstanceFsEntry,
		image string,
		net []*model.ContainerInstanceNetEntry,
		workDir string,
		containerId string,
		eventPublisher eventbus.EventPublisher,
		opGraphId string,
	) (err error)

	DeleteContainerIfExists(
		containerId string,
	)
}
