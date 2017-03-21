package docker

import (
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/opctl/opctl/util/containerprovider"
	"github.com/opctl/opctl/util/vfs"
	"github.com/opctl/opctl/util/vfs/os"
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

	containerProvider = _containerProvider{
		dockerClient: dockerClient,
		fs:           os.New(),
		runtime:      vruntime.New(),
	}

	return

}

type _containerProvider struct {
	dockerClient dockerClient
	fs           vfs.Vfs
	runtime      vruntime.Vruntime
}
