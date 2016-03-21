package osfilesys

import "io/ioutil"

type savePipelineFileUcExecuter interface {
  Execute(
  pipelineName string,
  data []byte,
  ) (err error)
}

func newSavePipelineFileUcExecuter(
relPathToPipelineFileFactory relPathToPipelineFileFactory,
) savePipelineFileUcExecuter {

  return &savePipelineFileUcExecuterImpl{
    relPathToPipelineFileFactory:relPathToPipelineFileFactory,
  }

}

type savePipelineFileUcExecuterImpl struct {
  relPathToPipelineFileFactory relPathToPipelineFileFactory
}

func (uc savePipelineFileUcExecuterImpl)  Execute(
pipelineName string,
data []byte,
) (err error) {

  var relPathToPipelineFile string
  relPathToPipelineFile, err = uc.relPathToPipelineFileFactory.Construct(pipelineName)
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