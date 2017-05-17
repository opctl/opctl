package docker

import (
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
	"io"
)

func (this _containerProvider) streamContainerStdOut(
	ctx context.Context,
	containerId string,
	writeCloser io.WriteCloser,
) error {

	readCloser, err := this.dockerClient.ContainerLogs(
		ctx,
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
	readCloser.Close()
	return err
}
