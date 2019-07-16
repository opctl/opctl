package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	dockerClientPkg "github.com/docker/docker/client"
	"io"
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
	containerID string,
	dst io.Writer,
) error {
	src, err := cses.dockerClient.ContainerLogs(
		ctx,
		containerID,
		types.ContainerLogsOptions{
			Follow:     true,
			ShowStderr: true,
		},
	)
	if nil != err {
		return err
	}

	_, err = io.Copy(dst, src)
	src.Close()
	return err
}
