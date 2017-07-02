package containerprovider

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ ContainerProvider

import (
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"io"
)

type ContainerProvider interface {
	DeleteContainerIfExists(
		containerId string,
	) error

	// RunContainer creates, starts, and waits on a container. ExitCode &/Or an error will be returned
	RunContainer(
		req *model.DCGContainerCall,
		eventPublisher pubsub.EventPublisher,
		stdout io.WriteCloser,
		stderr io.WriteCloser,
	) (*int64, error)
}
