package osfilesys

import "os"

type createPipelineDirUcExecuter interface {
  Execute(
  pipelineName string,
  ) (err error)
}

func newCreatePipelineDirUcExecuter(
relPathToPipelineDirFactory relPathToPipelineDirFactory,
) createPipelineDirUcExecuter {

  return &createPipelineDirUcExecuterImpl{
    relPathToPipelineDirFactory:relPathToPipelineDirFactory,
  }

}

type createPipelineDirUcExecuterImpl struct {
  relPathToPipelineDirFactory relPathToPipelineDirFactory
}

func (uc createPipelineDirUcExecuterImpl)  Execute(
pipelineName string,
) (err error) {

  var relPathToPipelineDir string
  relPathToPipelineDir, err = uc.relPathToPipelineDirFactory.Construct(pipelineName)
  if (nil != err) {
    return
  }

  err = os.MkdirAll(relPathToPipelineDir, 0777)

  return

}