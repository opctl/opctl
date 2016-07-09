package dockercompose

//go:generate counterfeiter -o ./fakeCompositionRoot.go --fake-name fakeCompositionRoot ./ compositionRoot

import (
  dockerEngine "github.com/docker/engine-api/client"
)

type compositionRoot interface {
  RunOpUseCase() runOpUseCase
  KillOpRunUseCase() killOpRunUseCase
}

func newCompositionRoot(
) (compositionRoot compositionRoot, err error) {

  filesys := _filesys{}

  dockerEngine, err := dockerEngine.NewEnvClient()
  if (nil != err) {
    return
  }

  opRunExitCodeReader := newOpRunExitCodeReader(dockerEngine)

  opRunResourceFlusher := newOpRunResourceFlusher()

  compositionRoot = &_compositionRoot{
    runOpUseCase: newRunOpUseCase(opRunExitCodeReader, opRunResourceFlusher),
    killOpRunUseCase: newKillOpRunUseCase(opRunResourceFlusher, filesys),
  }

  return
}

type _compositionRoot struct {
  runOpUseCase     runOpUseCase
  killOpRunUseCase killOpRunUseCase
}

func (this _compositionRoot) RunOpUseCase() runOpUseCase {
  return this.runOpUseCase
}

func (this _compositionRoot) KillOpRunUseCase() killOpRunUseCase {
  return this.killOpRunUseCase
}
