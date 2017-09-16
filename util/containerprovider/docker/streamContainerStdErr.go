package docker

import (
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
	"io"
)

func (ctp _containerProvider) streamContainerStdErr(
	containerId string,
	writeCloser io.WriteCloser,
) error {

	readCloser, err := ctp.dockerClient.ContainerLogs(
		context.TODO(),
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
	readCloser.Close()
	return err
}
