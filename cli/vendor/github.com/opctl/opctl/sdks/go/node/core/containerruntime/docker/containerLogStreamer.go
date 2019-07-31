package docker

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakeContainerLogStreamer.go --fake-name fakeContainerLogStreamer ./ containerLogStreamer

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
