package models

type PipelineStageRunView interface {
  StartedAtEpochTime() int64
  EndedAtEpochTime()   int64
  ExitCode()           int
}