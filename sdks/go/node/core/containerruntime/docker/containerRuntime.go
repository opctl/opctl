package docker

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o internal/fakes/commonAPIClient.go github.com/docker/docker/client.CommonAPIClient

import (
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
	"golang.org/x/net/context"
)

func New() (
	containerRuntime containerruntime.ContainerRuntime,
	err error,
) {

	dockerClient, err := dockerClientPkg.NewClientWithOpts(dockerClientPkg.FromEnv)
	if nil != err {
		return
	}

	// degrade client version to version of server
	dockerClient.NegotiateAPIVersion(context.TODO())

	rc, err := newRunContainer(dockerClient)
	if nil != err {
		return
	}

	objectUnderConstruction := _containerRuntime{
		runContainer: rc,
		dockerClient: dockerClient,
		os:           ios.New(),
	}
	containerRuntime = objectUnderConstruction

	// ensure user defined network exists to allow inter container resolution via name
	// @TODO: remove when socket outputs supported
	err = objectUnderConstruction.EnsureNetworkExists(dockerNetworkName)

	return
}

type _containerRuntime struct {
	runContainer
	dockerClient dockerClientPkg.CommonAPIClient
	os           ios.IOS
}

const dockerNetworkName = "opctl"
