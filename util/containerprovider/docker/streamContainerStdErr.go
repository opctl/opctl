package docker

import (
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
	"io"
)

func (this _containerProvider) streamContainerStdErr(
	containerId string,
	writeCloser io.WriteCloser,
) error {

	readCloser, err := this.dockerClient.ContainerLogs(
		context.Background(),
		containerId,
		types.ContainerLogsOptions{
			Follow:     true,
			ShowStderr: true,
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
