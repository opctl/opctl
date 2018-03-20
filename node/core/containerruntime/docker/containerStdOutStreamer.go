package docker

import (
	"github.com/docker/docker/api/types"
	dockerClientPkg "github.com/docker/docker/client"
	"golang.org/x/net/context"
	"io"
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
	containerID string,
	dst io.Writer,
) error {

	src, err := ctp.dockerClient.ContainerLogs(
		context.TODO(),
		containerID,
		types.ContainerLogsOptions{
			Follow:     true,
			ShowStdout: true,
		},
	)
	if nil != err {
		return err
	}

	_, err = io.Copy(dst, src)
	src.Close()
	return err
}
