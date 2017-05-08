package docker

import (
	dockerClientPkg "github.com/docker/docker/client"
	iclient "github.com/golang-interfaces/github.com-moby-moby/client"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/util/containerprovider"
	"github.com/opctl/opctl/util/vruntime"
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
	pingInfo, err := dockerClient.Ping(context.Background())
	if nil != err {
		return
	}
	dockerClient.UpdateClientVersion(pingInfo.APIVersion)

	objectUnderConstruction := _containerProvider{
		dockerClient: dockerClient,
		os:           ios.New(),
		runtime:      vruntime.New(),
	}
	containerProvider = objectUnderConstruction

	// ensure user defined network exists to allow inter container resolution via name
	// @TODO: remove when socket outputs supported
	err = objectUnderConstruction.EnsureNetworkExists(dockerNetworkName)

	return
}

type _containerProvider struct {
	dockerClient iclient.Client
	os           ios.IOS
	runtime      vruntime.Vruntime
}

const dockerNetworkName = "opctl"
