package os

import "io/ioutil"

type savePipelineFileUseCase interface {
  Execute(
  pipelineName string,
  data []byte,
  ) (err error)
}

func newSavePipelineFileUseCase(
relPathToPipelineFileFactory relPathToPipelineFileFactory,
) savePipelineFileUseCase {

  return &_savePipelineFileUseCase{
    relPathToPipelineFileFactory:relPathToPipelineFileFactory,
  }

}

type _savePipelineFileUseCase struct {
  relPathToPipelineFileFactory relPathToPipelineFileFactory
}

func (this _savePipelineFileUseCase)  Execute(
pipelineName string,
data []byte,
) (err error) {

  var relPathToPipelineFile string
  relPathToPipelineFile, err = this.relPathToPipelineFileFactory.Construct(pipelineName)
  if (nil != err) {
    return
  }

  err = ioutil.WriteFile(
    relPathToPipelineFile,
    data,
    0777,
  )

  return

}