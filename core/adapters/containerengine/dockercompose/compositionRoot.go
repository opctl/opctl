package dockercompose

import (
  dockerEngine "github.com/docker/engine-api/client"
)

type compositionRoot interface {
  InitOperationUseCase() initOperationUseCase
  RunOperationUseCase() runOperationUseCase
}

func newCompositionRoot(
) (compositionRoot compositionRoot, err error) {

  fs := filesystemImpl{}
  yml := _yamlCodec{}

  dockerEngine, err := dockerEngine.NewEnvClient()
  if (nil != err) {
    return
  }

  operationRunExitCodeReader := newOperationRunExitCodeReader(fs, dockerEngine)

  operationRunResourceFlusher := newOperationRunResourceFlusher(fs)

  compositionRoot = &_compositionRoot{
    initOperationUseCase: newInitOperationUseCase(fs, yml),
    runOperationUseCase: newRunOperationUseCase(fs, operationRunExitCodeReader, operationRunResourceFlusher),
  }

  return
}

type _compositionRoot struct {
  initOperationUseCase initOperationUseCase
  runOperationUseCase  runOperationUseCase
}

func (this _compositionRoot) InitOperationUseCase() initOperationUseCase {
  return this.initOperationUseCase
}

func (this _compositionRoot) RunOperationUseCase() runOperationUseCase {
  return this.runOperationUseCase
}
