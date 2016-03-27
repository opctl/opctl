package models

func NewPipelineStageRunView(
name               string,
stageType          string,
startedAtEpochTime int64,
endedAtEpochTime   int64,
exitCode           int,
stages             []*PipelineStageRunView,
) *PipelineStageRunView {

  return &PipelineStageRunView{
    Name:name,
    StageType:stageType,
    StartedAtEpochTime:startedAtEpochTime,
    EndedAtEpochTime:endedAtEpochTime,
    ExitCode:exitCode,
    Stages:stages,
  }

}

type PipelineStageRunView struct {
  Name               string
  StageType          string
  StartedAtEpochTime int64
  EndedAtEpochTime   int64
  ExitCode           int
  Stages             []*PipelineStageRunView `json:",omitempty"`
}