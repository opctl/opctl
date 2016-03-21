package dockercompose

import (
  dockerEngine "github.com/docker/engine-api/client"
)

type compositionRoot interface {
  InitDevOpUcExecuter() initDevOpUcExecuter
  RunDevOpUcExecuter() runDevOpUcExecuter
}

func newCompositionRoot(
) (compositionRoot compositionRoot, err error) {

  fs := filesystemImpl{}
  yml := yamlCodecImpl{}

  dockerEngine, err := dockerEngine.NewEnvClient()
  if (nil != err) {
    return
  }

  ecr := newDevOpExitCodeReader(fs, dockerEngine)

  rf := newDevOpResourceFlusher(fs)

  compositionRoot = &_compositionRoot{
    initDevOpUcExecuter: newInitDevOpUcExecuter(fs, yml),
    runDevOpUcExecuter: newRunDevOpUcExecuter(fs, ecr, rf),
  }

  return
}

type _compositionRoot struct {
  initDevOpUcExecuter initDevOpUcExecuter
  runDevOpUcExecuter  runDevOpUcExecuter
}

func (_compositionRoot _compositionRoot) InitDevOpUcExecuter() initDevOpUcExecuter {
  return _compositionRoot.initDevOpUcExecuter
}

func (_compositionRoot _compositionRoot) RunDevOpUcExecuter() runDevOpUcExecuter {
  return _compositionRoot.runDevOpUcExecuter
}
