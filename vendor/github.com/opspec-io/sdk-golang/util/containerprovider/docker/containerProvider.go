package docker

//go:generate counterfeiter -o ./fakeDockerClient.go --fake-name fakeDockerClient /go/src/github.com/opspec-io/sdk-golang/vendor/github.com/docker/docker/client/interface.go CommonAPIClient

import (
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/util/containerprovider"
	"golang.org/x/net/context"
)

func New() (
	containerProvider containerprovider.ContainerProvider,
	err error,
) {

	dockerClient, err := dockerClientPkg.NewEnvClient()
	if nil != err {
		return
	}

	// degrade client version to version of server
	dockerClient.NegotiateAPIVersion(context.TODO())

	objectUnderConstruction := _containerProvider{
		runContainer: newRunContainer(dockerClient),
		dockerClient: dockerClient,
		os:           ios.New(),
	}
	containerProvider = objectUnderConstruction

	// ensure user defined network exists to allow inter container resolution via name
	// @TODO: remove when socket outputs supported
	err = objectUnderConstruction.EnsureNetworkExists(dockerNetworkName)

	return
}

type _containerProvider struct {
	runContainer
	dockerClient dockerClientPkg.CommonAPIClient
	os           ios.IOS
}

const dockerNetworkName = "opctl"
