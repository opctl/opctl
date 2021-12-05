package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	dockerClientPkg "github.com/docker/docker/client"
)

func newContainerStdErrStreamer(
	dockerClient dockerClientPkg.CommonAPIClient,
) containerLogStreamer {
	return _containerStdErrStreamer{
		dockerClient,
	}
}

type _containerStdErrStreamer struct {
	dockerClient dockerClientPkg.CommonAPIClient
}

func (cses _containerStdErrStreamer) Stream(
	ctx context.Context,
	containerName string,
	dst io.Writer,
) error {
	src, err := cses.dockerClient.ContainerLogs(
		ctx,
		containerName,
		types.ContainerLogsOptions{
			Follow:     true,
			ShowStderr: true,
		},
	)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, src)
	src.Close()
	return err
}
