package models

type PipelineRunView struct {
  Id                *string
  PipelineName      string
  Stages            []*PipelineStageRunView
  StartedAtUnixTime int64
  EndedAtUnixTime   int64
  ExitCode          int
}