package docker

import (
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/opctl/opctl/util/containerprovider"
	"github.com/opctl/opctl/util/vruntime"
	"github.com/virtual-go/vfs"
	"github.com/virtual-go/vfs/osfs"
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
		fs:           osfs.New(),
		runtime:      vruntime.New(),
	}
	containerProvider = objectUnderConstruction

	// ensure user defined network exists to allow inter container resolution via name
	// @TODO: remove when socket outputs supported
	err = objectUnderConstruction.EnsureNetworkExists("opctl")

	return
}

type _containerProvider struct {
	dockerClient dockerClient
	fs           vfs.Vfs
	runtime      vruntime.Vruntime
}
