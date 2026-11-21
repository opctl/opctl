package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/container"
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

//counterfeiter:generate -o internal/fakes/containerLogStreamer.go . containerLogStreamer
type containerLogStreamer interface {
	Stream(
		ctx context.Context,
		containerName string,
		stdout io.Writer,
		stderr io.Writer,
	) error
}

func newContainerLogStreamer(
	dockerClient dockerClientPkg.CommonAPIClient,
) containerLogStreamer {
	return _containerLogStreamer{
		dockerClient,
	}
}

type _containerLogStreamer struct {
	dockerClient dockerClientPkg.CommonAPIClient
}

// Stream follows stdout and stderr through one log stream so lines keep the
// order the container wrote them in.
func (cls _containerLogStreamer) Stream(
	ctx context.Context,
	containerName string,
	stdout io.Writer,
	stderr io.Writer,
) error {
	src, err := cls.dockerClient.ContainerLogs(
		ctx,
		containerName,
		container.LogsOptions{
			Follow:     true,
			ShowStdout: true,
			ShowStderr: true,
		},
	)
	if err != nil {
		return err
	}
	defer src.Close()

	// without a TTY docker multiplexes stdout/stderr frames
	_, err = stdcopy.StdCopy(stdout, stderr, src)
	return err
}
