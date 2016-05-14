package dockercompose

//go:generate counterfeiter -o ./fakeCompositionRoot.go --fake-name fakeCompositionRoot ./ compositionRoot

import (
  dockerEngine "github.com/docker/engine-api/client"
)

type compositionRoot interface {
  InitOpUseCase() initOpUseCase
  RunOpUseCase() runOpUseCase
  KillOpRunUseCase() killOpRunUseCase
}

func newCompositionRoot(
) (compositionRoot compositionRoot, err error) {

  fs := filesystemImpl{}
  yml := _yamlCodec{}

  dockerEngine, err := dockerEngine.NewEnvClient()
  if (nil != err) {
    return
  }

  opRunExitCodeReader := newOpRunExitCodeReader(dockerEngine)

  opRunResourceFlusher := newOpRunResourceFlusher()

  compositionRoot = &_compositionRoot{
    initOpUseCase: newInitOpUseCase(fs, yml),
    runOpUseCase: newRunOpUseCase(opRunExitCodeReader, opRunResourceFlusher),
    killOpRunUseCase: newKillOpRunUseCase(opRunResourceFlusher),
  }

  return
}

type _compositionRoot struct {
  initOpUseCase    initOpUseCase
  runOpUseCase     runOpUseCase
  killOpRunUseCase killOpRunUseCase
}

func (this _compositionRoot) InitOpUseCase() initOpUseCase {
  return this.initOpUseCase
}

func (this _compositionRoot) RunOpUseCase() runOpUseCase {
  return this.runOpUseCase
}

func (this _compositionRoot) KillOpRunUseCase() killOpRunUseCase {
  return this.killOpRunUseCase
}
