package dockercompose

import (
  dockerEngine "github.com/docker/engine-api/client"
)

type compositionRoot interface {
  InitDevOpUseCase() initDevOpUseCase
  RunDevOpUseCase() runDevOpUseCase
}

func newCompositionRoot(
) (compositionRoot compositionRoot, err error) {

  fs := filesystemImpl{}
  yml := _yamlCodec{}

  dockerEngine, err := dockerEngine.NewEnvClient()
  if (nil != err) {
    return
  }

  devOpRunExitCodeReader := newDevOpRunExitCodeReader(fs, dockerEngine)

  devOpRunResourceFlusher := newDevOpRunResourceFlusher(fs)

  compositionRoot = &_compositionRoot{
    initDevOpUseCase: newInitDevOpUseCase(fs, yml),
    runDevOpUseCase: newRunDevOpUseCase(fs, devOpRunExitCodeReader, devOpRunResourceFlusher),
  }

  return
}

type _compositionRoot struct {
  initDevOpUseCase initDevOpUseCase
  runDevOpUseCase  runDevOpUseCase
}

func (this _compositionRoot) InitDevOpUseCase() initDevOpUseCase {
  return this.initDevOpUseCase
}

func (this _compositionRoot) RunDevOpUseCase() runDevOpUseCase {
  return this.runDevOpUseCase
}
