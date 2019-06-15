package docker

//go:generate counterfeiter -o ./fakeContainerLogStreamer.go --fake-name fakeContainerLogStreamer ./ containerLogStreamer

import (
	"context"
	"io"
)

type containerLogStreamer interface {
	Stream(
		ctx context.Context,
		containerID string,
		dst io.Writer,
	) error
}
