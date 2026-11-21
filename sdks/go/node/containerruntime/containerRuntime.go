// Package containerruntime defines an interface abstracting container runtime interactions.
// A fake implementation is included to allow faking said interactions.
package containerruntime

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"io"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

// ContainerRuntime defines the interface container runtimes must implement to be supported by opctl
//
//counterfeiter:generate -o fakes/containerRuntime.go . ContainerRuntime
type ContainerRuntime interface {
	// Delete deletes opctl managed resources from the container runtime
	Delete(
		ctx context.Context,
	) error

	DeleteContainerIfExists(
		ctx context.Context,
		containerID string,
	) error

	// Kill stops/kills opctl managed resources within the container runtime
	Kill(
		ctx context.Context,
	) error

	// RunContainer creates, starts, and waits on a container.
	// stdin and resize are nil unless interactive input was requested and is supported.
	RunContainer(
		ctx context.Context,
		req *model.ContainerCall,
		// @TODO: get rid of in combination with eventPublisher
		rootCallID string,
		// @TODO: get rid of this; just use stdout/stderr
		eventPublisher pubsub.EventPublisher,
		stdout io.WriteCloser,
		stderr io.WriteCloser,
		stdin io.Reader,
		resize <-chan model.TermSize,
	) (*int64, error)
}
