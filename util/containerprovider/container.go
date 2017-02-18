package containerprovider

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ ContainerProvider

import (
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

type ContainerProvider interface {
	InspectContainerIfExists(
		containerId string,
	) (container *model.DcgContainerCall, err error)

	DeleteContainerIfExists(
		containerId string,
	)

	RunContainer(
		req *RunContainerReq,
		eventPublisher pubsub.EventPublisher,
	) (err error)
}
