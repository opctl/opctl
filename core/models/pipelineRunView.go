package models

type PipelineRunView struct {
  PipelineName       string
  Stages             []*PipelineStageRunView
  StartedAtUnixTime int64
  EndedAtUnixTime   int64
  ExitCode           int
}

func newPipelineRunView(
pipelineName       string,
stages          []*PipelineStageRunView,
startedAtUnixTime int64,
endedAtUnixTime   int64,
exitCode           int,
) PipelineRunView {

  return PipelineRunView{
    PipelineName:pipelineName,
    Stages:stages,
    StartedAtUnixTime:startedAtUnixTime,
    EndedAtUnixTime:endedAtUnixTime,
    ExitCode:exitCode,
  }

}