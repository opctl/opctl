package osfilesys

import "os"

type createPipelineDirUseCase interface {
  Execute(
  pipelineName string,
  ) (err error)
}

func newCreatePipelineDirUseCase(
relPathToPipelineDirFactory relPathToPipelineDirFactory,
) createPipelineDirUseCase {

  return &_createPipelineDirUseCase{
    relPathToPipelineDirFactory:relPathToPipelineDirFactory,
  }

}

type _createPipelineDirUseCase struct {
  relPathToPipelineDirFactory relPathToPipelineDirFactory
}

func (this _createPipelineDirUseCase)  Execute(
pipelineName string,
) (err error) {

  var relPathToPipelineDir string
  relPathToPipelineDir, err = this.relPathToPipelineDirFactory.Construct(pipelineName)
  if (nil != err) {
    return
  }

  err = os.MkdirAll(relPathToPipelineDir, 0777)

  return

}