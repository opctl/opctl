package docker

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o internal/fakes/commonAPIClient.go github.com/docker/docker/client.CommonAPIClient

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
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
	containers, err := cr.dockerClient.ContainerList(
		ctx,
		container.ListOptions{
			Filters: filters.NewArgs(
				filters.KeyValuePair{
					Key:   "name",
					Value: containerNamePrefix,
				},
				filters.KeyValuePair{
					Key:   "network",
					Value: networkName,
				},
			),
		},
	)
	if err != nil {
		return err
	}

	errGroup, egCtx := errgroup.WithContext(ctx)
	for _, container := range containers {
		for _, containerName := range container.Names {
			containerName := containerName
			errGroup.Go(func() error {
				slashPrefix := fmt.Sprintf("/%s", containerNamePrefix)
				// check if containerName is a conventional opctl container name
				if strings.HasPrefix(containerName, slashPrefix) {
					return cr.DeleteContainerIfExists(
						egCtx,
						// convert containerName to opctl container id as required by cr.DeleteContainerIfExists
						strings.Replace(containerName, slashPrefix, "", 1),
					)
				}
				return nil
			})
		}
	}

	err = errGroup.Wait()
	if err != nil {
		return err
	}

	return ensureNetworkDetached(ctx, cr.dockerClient)

}

func (cr _containerRuntime) Kill(
	ctx context.Context,
) error {
	return cr.Delete(ctx)
}

const containerNamePrefix = "opctl_"
const networkName = "opctl"
