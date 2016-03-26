package os

import "path"

type relPathToPipelineFileFactory interface {
  Construct(
  pipelineName string,
  ) (relPathToPipelineFile string, err error)
}

func newRelPathToPipelineFileFactory(
relPathToPipelineDirFactory relPathToPipelineDirFactory,
) relPathToPipelineFileFactory {

  return &_relPathToPipelineFileFactory{
    relPathToPipelineDirFactory:relPathToPipelineDirFactory,
  }

}

type _relPathToPipelineFileFactory struct {
  relPathToPipelineDirFactory relPathToPipelineDirFactory
}

func (this _relPathToPipelineFileFactory) Construct(
pipelineName string,
) (relPathToPipelineFile string, err error) {

  var relPathToPipelineDir string
  relPathToPipelineDir, err = this.relPathToPipelineDirFactory.Construct(pipelineName)
  if (nil != err) {
    return
  }

  relPathToPipelineFile = path.Join(relPathToPipelineDir, "pipeline.yml")

  return

}