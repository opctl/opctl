package core

import (
"github.com/dev-op-spec/engine/core/models"
"github.com/dev-op-spec/engine/core/ports"
)

type runDevOpUseCase interface {
  Execute(
  devOpName string,
  ) (devOpRun models.DevOpRunView, err error)
}

func newRunDevOpUseCase(
ce ports.ContainerEngine,
) runDevOpUseCase {

  return &_runDevOpUseCase{
    ce:ce,
  }

}

type _runDevOpUseCase struct {
  ce ports.ContainerEngine
}

func (this _runDevOpUseCase) Execute(
devOpName string,
) (devOpRun models.DevOpRunView, err error) {

  return this.ce.RunDevOp(devOpName)

}
