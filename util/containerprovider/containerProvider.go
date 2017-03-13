package containerprovider

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ ContainerProvider

import (
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

type ContainerProvider interface {
	CreateNetwork(
		networkId string,
	) (err error)

	DeleteNetworkIfExists(
		networkId string,
	) (err error)

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
