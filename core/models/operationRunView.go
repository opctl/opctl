package models

type OperationRunView struct {
  Id                *string
  OperationName     string
  SubOperations     []*SubOperationRunView `json:",omitempty"`
  StartedAtUnixTime int64
  EndedAtUnixTime   int64
  ExitCode          int
}