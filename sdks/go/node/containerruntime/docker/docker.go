package docker

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o internal/fakes/commonAPIClient.go github.com/docker/docker/client.CommonAPIClient

import (
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"golang.org/x/net/context"
)

func New(
	ctx context.Context,
	host string,
) (containerruntime.ContainerRuntime, error) {

	dockerClient, err := dockerClientPkg.NewClientWithOpts(dockerClientPkg.FromEnv, dockerClientPkg.WithHost(host))
	if err != nil {
		return nil, err
	}

	// degrade client version to version of server
	dockerClient.NegotiateAPIVersion(ctx)

	rc, err := newRunContainer(ctx, dockerClient)
	if err != nil {
		return nil, err
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

func (cr _containerRuntime) Delete(
	ctx context.Context,
) error {
	// for now this is a no-op
	return nil
}

const dockerNetworkName = "opctl"
