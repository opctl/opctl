package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type runDevOpUseCase interface {
  Execute(
  req models.RunDevOpReq,
  ) (devOpRun models.DevOpRunView, err error)
}

func newRunDevOpUseCase(
containerEngine ports.ContainerEngine,
pathToDevOpDirFactory pathToDevOpDirFactory,
) runDevOpUseCase {

  return &_runDevOpUseCase{
    containerEngine:containerEngine,
    pathToDevOpDirFactory:pathToDevOpDirFactory,
  }

}

type _runDevOpUseCase struct {
  containerEngine       ports.ContainerEngine
  pathToDevOpDirFactory pathToDevOpDirFactory
}

func (this _runDevOpUseCase) Execute(
req models.RunDevOpReq,
) (devOpRun models.DevOpRunView, err error) {

  pathToDevOpDir := this.pathToDevOpDirFactory.Construct(
    req.ProjectUrl,
    req.DevOpName,
  )

  devOpRun, err = this.containerEngine.RunDevOp(pathToDevOpDir)

  return
}
