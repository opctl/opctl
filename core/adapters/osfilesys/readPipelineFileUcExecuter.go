package osfilesys

import "io/ioutil"

type readPipelineFileUcExecuter interface {
  Execute(
  pipelineName string,
  ) (pipelineFile []byte, err error)
}

func newReadPipelineFileUcExecuter(
relPathToPipelineFileFactory relPathToPipelineDirFactory,
) readPipelineFileUcExecuter {

  return &readPipelineFileUcExecuterImpl{
    relPathToPipelineFileFactory:relPathToPipelineFileFactory,
  }

}

type readPipelineFileUcExecuterImpl struct {
  relPathToPipelineFileFactory relPathToPipelineDirFactory
}

func (uc readPipelineFileUcExecuterImpl) Execute(
pipelineName string,
) (pipelineFile []byte, err error) {

  relativePathToPipelineFile, err := uc.relPathToPipelineFileFactory.Construct(pipelineName)
  if (nil != err) {
    return
  }

  pipelineFile, err = ioutil.ReadFile(relativePathToPipelineFile)

  return
}
