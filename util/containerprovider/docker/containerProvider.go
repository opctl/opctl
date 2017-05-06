package docker

import (
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/golang-interfaces/vos"
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
		os:           vos.New(),
		runtime:      vruntime.New(),
	}
	containerProvider = objectUnderConstruction

	// ensure user defined network exists to allow inter container resolution via name
	// @TODO: remove when socket outputs supported
	err = objectUnderConstruction.EnsureNetworkExists(dockerNetworkName)

	return
}

type _containerProvider struct {
	dockerClient dockerClient
	os           vos.VOS
	runtime      vruntime.Vruntime
}

const dockerNetworkName = "opctl"
