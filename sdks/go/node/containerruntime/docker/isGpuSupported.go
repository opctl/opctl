package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
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
	containerName := getContainerName(fmt.Sprintf("gpu-check-%s", uuid.NewV4().String()))

	defer dockerClient.ContainerRemove(
		context.Background(),
		containerName,
		types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		},
	)

	imageRef := "alpine"

	err := pullImage(
		context.Background(),
		&model.ContainerCall{
			Image: &model.ContainerCallImage{
				Ref:       &imageRef,
				PullCreds: imagePullCreds,
			},
		},
		dockerClient,
		"",
		noOpEventPublisher{},
	)

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
		types.ContainerStartOptions{},
	)
	if startErr != nil {
		return false, nil
	}

	return true, nil
}
