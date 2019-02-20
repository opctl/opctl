package docker

import (
	"github.com/docker/docker/api/types"
	dockerClientPkg "github.com/docker/docker/client"
	"golang.org/x/net/context"
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
	containerID string,
	dst io.Writer,
) error {
	src, err := cses.dockerClient.ContainerLogs(
		context.TODO(),
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
