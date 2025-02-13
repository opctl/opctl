package docker

import (
	"context"
	"fmt"
	"runtime"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/satori/go.uuid"
)

type noOpEventPublisher struct{}

func (ep noOpEventPublisher) Publish(
	event model.Event,
) {
}

func isGpuSupported(
	ctx context.Context,
	dockerClient client.CommonAPIClient,
	imagePullCreds *model.Creds,
) (bool, error) {
	if runtime.GOOS != "linux" {
		// GPU passthrough only available when host is linux
		return false, nil
	}

	containerName := getContainerName(fmt.Sprintf("gpu-check-%s", uuid.NewV4().String()))

	defer dockerClient.ContainerRemove(
		context.Background(),
		containerName,
		container.RemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		},
	)

	imageRef := "ghcr.io/linuxcontainers/alpine"

	if err := pullImage(
		context.Background(),
		&model.ContainerCall{
			Image: &model.ContainerCallImage{
				Ref: &imageRef,
			},
		},
		dockerClient,
		"",
		noOpEventPublisher{},
	); err != nil {
		return false, err
	}

	createResponse, err := dockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageRef,
			Cmd:   []string{"echo"},
		},
		&container.HostConfig{
			AutoRemove: true,
			Resources: container.Resources{
				DeviceRequests: []container.DeviceRequest{
					{
						Capabilities: [][]string{{"gpu"}},
						Count:        -1,
					},
				},
			},
		},
		&network.NetworkingConfig{},
		nil,
		containerName,
	)
	if nil != err {
		return false, err
	}

	startErr := dockerClient.ContainerStart(
		ctx,
		createResponse.ID,
		container.StartOptions{},
	)
	if startErr != nil {
		return false, nil
	}

	return true, nil
}
