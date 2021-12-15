package docker

import (
	"context"
	"io"
)

//counterfeiter:generate -o internal/fakes/containerLogStreamer.go . containerLogStreamer
type containerLogStreamer interface {
	Stream(
		ctx context.Context,
		containerName string,
		dst io.Writer,
	) error
}
