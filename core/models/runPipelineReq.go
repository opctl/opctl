package models

func NewRunPipelineReq(
pathToProjectRootDir string,
pipelineName string,
) *RunPipelineReq {

  return &RunPipelineReq{
    PathToProjectRootDir:pathToProjectRootDir,
    PipelineName :pipelineName,
  }

}

type RunPipelineReq struct {
  PathToProjectRootDir string
  PipelineName         string
}
