package docker

import (
  dockerClientPkg "github.com/docker/docker/client"
  "github.com/opspec-io/engine/pkg/containerengine"
)

func New(
) (
containerEngine containerengine.ContainerEngine,
err error,
) {

  dockerClient, err := dockerClientPkg.NewEnvClient()
  if (nil != err) {
    return
  }

  containerEngine = _containerEngine{
    containerExitCodeReader :newContainerExitCodeReader(dockerClient),
  }

  return

}

type _containerEngine struct {
  containerExitCodeReader containerExitCodeReader
}
