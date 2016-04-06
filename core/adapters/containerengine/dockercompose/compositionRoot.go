package dockercompose

import (
  dockerEngine "github.com/docker/engine-api/client"
)

type compositionRoot interface {
  InitOpUseCase() initOpUseCase
  RunOpUseCase() runOpUseCase
}

func newCompositionRoot(
) (compositionRoot compositionRoot, err error) {

  fs := filesystemImpl{}
  yml := _yamlCodec{}

  dockerEngine, err := dockerEngine.NewEnvClient()
  if (nil != err) {
    return
  }

  opRunExitCodeReader := newOpRunExitCodeReader(fs, dockerEngine)

  opRunResourceFlusher := newOpRunResourceFlusher(fs)

  compositionRoot = &_compositionRoot{
    initOpUseCase: newInitOpUseCase(fs, yml),
    runOpUseCase: newRunOpUseCase(opRunExitCodeReader, opRunResourceFlusher),
  }

  return
}

type _compositionRoot struct {
  initOpUseCase initOpUseCase
  runOpUseCase  runOpUseCase
}

func (this _compositionRoot) InitOpUseCase() initOpUseCase {
  return this.initOpUseCase
}

func (this _compositionRoot) RunOpUseCase() runOpUseCase {
  return this.runOpUseCase
}
