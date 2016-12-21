package docker

import (
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/opspec-io/engine/pkg/containerengine"
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
	}

	return

}

type _containerEngine struct {
	dockerEngine dockerEngine
}
