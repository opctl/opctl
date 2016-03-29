package models

func NewPipelineStageRunView(
name               string,
stageType          string,
startedAtUnixTime int64,
endedAtUnixTime   int64,
exitCode           int,
stages             []*PipelineStageRunView,
) *PipelineStageRunView {

  return &PipelineStageRunView{
    Name:name,
    StageType:stageType,
    StartedAtUnixTime:startedAtUnixTime,
    EndedAtUnixTime:endedAtUnixTime,
    ExitCode:exitCode,
    Stages:stages,
  }

}

type PipelineStageRunView struct {
  Name               string
  StageType          string
  StartedAtUnixTime int64
  EndedAtUnixTime   int64
  ExitCode           int
  Stages             []*PipelineStageRunView `json:",omitempty"`
}