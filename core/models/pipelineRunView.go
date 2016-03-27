package models

type PipelineRunView struct {
  PipelineName       string
  Stages             []*PipelineStageRunView
  StartedAtEpochTime int64
  EndedAtEpochTime   int64
  ExitCode           int
}

func newPipelineRunView(
pipelineName       string,
stages          []*PipelineStageRunView,
startedAtEpochTime int64,
endedAtEpochTime   int64,
exitCode           int,
) PipelineRunView {

  return PipelineRunView{
    PipelineName:pipelineName,
    Stages:stages,
    StartedAtEpochTime:startedAtEpochTime,
    EndedAtEpochTime:endedAtEpochTime,
    ExitCode:exitCode,
  }

}