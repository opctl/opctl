package docker

import (
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/opspec-io/opctl/pkg/containerengine"
	"github.com/opspec-io/opctl/util/fs"
	"github.com/opspec-io/opctl/util/fs/os"
	"golang.org/x/net/context"
)

func New() (
	containerEngine containerengine.ContainerEngine,
	err error,
) {

	dockerEngine, err := dockerClientPkg.NewEnvClient()
	if nil != err {
		return
	}

	// degrade client version to version of server
	pingInfo, err := dockerEngine.Ping(context.Background())
	if nil != err {
		return
	}
	dockerEngine.UpdateClientVersion(pingInfo.APIVersion)

	containerEngine = _containerEngine{
		dockerEngine: dockerEngine,
		fs:           os.New(),
	}

	return

}

type _containerEngine struct {
	dockerEngine dockerEngine
	fs           fs.Fs
}
