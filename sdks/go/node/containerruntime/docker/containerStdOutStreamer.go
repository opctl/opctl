package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	dockerClientPkg "github.com/docker/docker/client"
)

func newContainerStdOutStreamer(
	dockerClient dockerClientPkg.CommonAPIClient,
) containerLogStreamer {
	return _containerStdOutStreamer{
		dockerClient,
	}
}

type _containerStdOutStreamer struct {
	dockerClient dockerClientPkg.CommonAPIClient
}

func (ctp _containerStdOutStreamer) Stream(
	ctx context.Context,
	containerName string,
	dst io.Writer,
) error {

	src, err := ctp.dockerClient.ContainerLogs(
		ctx,
		containerName,
		types.ContainerLogsOptions{
			Follow:     true,
			ShowStdout: true,
		},
	)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, src)
	src.Close()
	return err
}
