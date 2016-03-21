package core

import (
"github.com/dev-op-spec/engine/core/models"
"github.com/dev-op-spec/engine/core/ports"
)

type runDevOpUcExecuter interface {
  Execute(
  devOpName string,
  ) (devOpRun models.DevOpRunView, err error)
}

func newRunDevOpUcExecuter(
ce ports.ContainerEngine,
) runDevOpUcExecuter {

  return &runDevOpUcExecuterImpl{
    ce:ce,
  }

}

type runDevOpUcExecuterImpl struct {
  ce ports.ContainerEngine
}

func (x runDevOpUcExecuterImpl) Execute(
devOpName string,
) (devOpRun models.DevOpRunView, err error) {

  return x.ce.RunDevOp(devOpName)

}
