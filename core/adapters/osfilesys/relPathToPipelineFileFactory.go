package osfilesys

import "path"

type relPathToPipelineFileFactory interface {
  Construct(
  pipelineName string,
  ) (relPathToPipelineFile string, err error)
}

func newRelPathToPipelineFileFactory(
relPathToPipelineDirFactory relPathToPipelineDirFactory,
) relPathToPipelineFileFactory {

  return &relPathToPipelineFileFactoryImpl{
    relPathToPipelineDirFactory:relPathToPipelineDirFactory,
  }

}

type relPathToPipelineFileFactoryImpl struct {
  relPathToPipelineDirFactory relPathToPipelineDirFactory
}

func (f relPathToPipelineFileFactoryImpl) Construct(
pipelineName string,
) (relPathToPipelineFile string, err error) {

  var relPathToPipelineDir string
  relPathToPipelineDir, err = f.relPathToPipelineDirFactory.Construct(pipelineName)
  if (nil != err) {
    return
  }

  relPathToPipelineFile = path.Join(relPathToPipelineDir, "pipeline.yml")

  return

}