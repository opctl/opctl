package docker

//go:generate counterfeiter -o ./fakeDockerClient.go --fake-name fakeDockerClient ./ dockerClient

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"golang.org/x/net/context"
)

type dockerClient interface {
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error)
	ContainerKill(ctx context.Context, container, signal string) error
	ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error)
	ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error
	ContainerStart(ctx context.Context, container string, options types.ContainerStartOptions) error
	ContainerWait(ctx context.Context, container string, condition container.WaitCondition) (<-chan container.ContainerWaitOKBody, <-chan error)
	ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error)
	NetworkCreate(ctx context.Context, name string, options types.NetworkCreate) (types.NetworkCreateResponse, error)
	NetworkInspect(ctx context.Context, networkID string, options types.NetworkInspectOptions) (types.NetworkResource, error)
}
