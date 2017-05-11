package docker

import (
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
	"io"
)

func (this _containerProvider) streamContainerStdOut(
	containerId string,
	writeCloser io.WriteCloser,
) error {

	readCloser, err := this.dockerClient.ContainerLogs(
		context.Background(),
		containerId,
		types.ContainerLogsOptions{
			Follow:     true,
			ShowStdout: true,
		},
	)
	if nil != err {
		return err
	}

	_, err = io.Copy(writeCloser, readCloser)
	writeCloser.Close()
	readCloser.Close()
	return err
}
