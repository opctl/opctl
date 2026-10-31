package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/container"
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
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
		container.LogsOptions{
			Follow:     true,
			ShowStdout: true,
		},
	)
	if err != nil {
		return err
	}

	if containerruntime.SeparateOutputFrom(ctx) {
		_, err = stdcopy.StdCopy(dst, io.Discard, src)
	} else {
		_, err = io.Copy(dst, src)
	}
	src.Close()
	return err
}
