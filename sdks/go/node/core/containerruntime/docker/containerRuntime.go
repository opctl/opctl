package docker

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o internal/fakes/commonAPIClient.go github.com/docker/docker/client.CommonAPIClient

import (
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
	"golang.org/x/net/context"
)

func New(ctx context.Context) (
	containerRuntime containerruntime.ContainerRuntime,
	err error,
) {

	dockerClient, err := dockerClientPkg.NewClientWithOpts(dockerClientPkg.FromEnv)
	if err != nil {
		return
	}

	// degrade client version to version of server
	dockerClient.NegotiateAPIVersion(ctx)

	rc, err := newRunContainer(ctx, dockerClient)
	if err != nil {
		return
	}

	return _containerRuntime{
		runContainer: rc,
		dockerClient: dockerClient,
	}, nil
}

type _containerRuntime struct {
	runContainer
	dockerClient dockerClientPkg.CommonAPIClient
}

const dockerNetworkName = "opctl"
