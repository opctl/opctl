package os

import "io/ioutil"

type readPipelineFileUseCase interface {
  Execute(
  pipelineName string,
  ) (pipelineFile []byte, err error)
}

func newReadPipelineFileUseCase(
relPathToPipelineFileFactory relPathToPipelineDirFactory,
) readPipelineFileUseCase {

  return &_readPipelineFileUseCase{
    relPathToPipelineFileFactory:relPathToPipelineFileFactory,
  }

}

type _readPipelineFileUseCase struct {
  relPathToPipelineFileFactory relPathToPipelineDirFactory
}

func (this _readPipelineFileUseCase) Execute(
pipelineName string,
) (pipelineFile []byte, err error) {

  relativePathToPipelineFile, err := this.relPathToPipelineFileFactory.Construct(pipelineName)
  if (nil != err) {
    return
  }

  pipelineFile, err = ioutil.ReadFile(relativePathToPipelineFile)

  return
}
