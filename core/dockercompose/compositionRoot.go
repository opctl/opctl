package dockercompose

//go:generate counterfeiter -o ./fakeCompositionRoot.go --fake-name fakeCompositionRoot ./ compositionRoot

import (
  dockerClientPkg "github.com/docker/docker/client"
)

type compositionRoot interface {
  StartContainerUseCase() startContainerUseCase
  EnsureContainerRemovedUseCase() ensureContainerRemovedUseCase
}

func newCompositionRoot(
) (compositionRoot compositionRoot, err error) {

  filesys := _filesys{}

  dockerClient, err := dockerClientPkg.NewEnvClient()
  if (nil != err) {
    return
  }

  containerExitCodeReader := newContainerExitCodeReader(dockerClient)

  containerRemover := newContainerRemover()

  compositionRoot = &_compositionRoot{
    startContainerUseCase: newStartOpRunUseCase(containerExitCodeReader, containerRemover),
    ensureContainerRemovedUseCase: newEnsureContainerRemovedUseCase(containerRemover, filesys),
  }

  return
}

type _compositionRoot struct {
  startContainerUseCase         startContainerUseCase
  ensureContainerRemovedUseCase ensureContainerRemovedUseCase
}

func (this _compositionRoot) StartOpRunUseCase() startContainerUseCase {
  return this.startContainerUseCase
}

func (this _compositionRoot) KillOpRunUseCase() ensureContainerRemovedUseCase {
  return this.ensureContainerRemovedUseCase
}
