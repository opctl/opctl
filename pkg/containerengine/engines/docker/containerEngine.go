package docker

import (
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/opspec-io/opctl/pkg/containerengine"
	"github.com/opspec-io/opctl/util/vfs"
	"github.com/opspec-io/opctl/util/vfs/os"
	"golang.org/x/net/context"
)

func New() (
	containerEngine containerengine.ContainerEngine,
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

	containerEngine = _containerEngine{
		dockerClient: dockerClient,
		vfs:          os.New(),
	}

	return

}

type _containerEngine struct {
	dockerClient dockerClient
	vfs          vfs.Vfs
}
