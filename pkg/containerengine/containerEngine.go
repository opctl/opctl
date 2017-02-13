package containerengine

//go:generate counterfeiter -o engines/fake/containerEngine.go --fake-name ContainerEngine ./ ContainerEngine

import (
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

type ContainerEngine interface {
	InspectContainerIfExists(
		containerId string,
	) (container *model.DcgContainerCall, err error)

	DeleteContainerIfExists(
		containerId string,
	)

	StartContainer(
		req *StartContainerReq,
		eventPublisher pubsub.EventPublisher,
	) (err error)
}
