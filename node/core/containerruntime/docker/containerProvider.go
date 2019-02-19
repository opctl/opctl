package docker

//go:generate counterfeiter -o ./fakeDockerClient.go --fake-name fakeDockerClient /go/src/github.com/opctl/sdk-golang/vendor/github.com/docker/docker/client/interface.go CommonAPIClient

import (
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/sdk-golang/node/core/containerruntime"
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

	objectUnderConstruction := _containerRuntime{
		runContainer: newRunContainer(dockerClient),
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
