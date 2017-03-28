package containerprovider

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ ContainerProvider

import (
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/model"
)

type ContainerProvider interface {
	DeleteContainerIfExists(
		containerId string,
	) (err error)

	NetworkContainer(
		networkId string,
		containerId string,
		containerAlias string,
	) (err error)

	RunContainer(
		req *model.DCGContainerCall,
		eventPublisher pubsub.EventPublisher,
	) (err error)
}
