package docker

//go:generate counterfeiter -o ./fakeContainerLogStreamer.go --fake-name fakeContainerLogStreamer ./ containerLogStreamer

import (
	"io"
)

type containerLogStreamer interface {
	Stream(
		containerID string,
		dst io.Writer,
	) error
}
