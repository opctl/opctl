package models

type PipelineStageRunView interface {
  StartedAtPosixTime() int64
  EndedAtPosixTime()   int64
  ExitCode()           int
}