package models

type OperationRunDetailedView struct {
  Id                *string `json:"id"`
  OperationUrl      *Url `json:"operationUrl"`
  SubOperations     []*OperationRunSummaryView `json:"subOperations,omitempty"`
  StartedAtUnixTime int64 `json:"startedAtUnixTime"`
  EndedAtUnixTime   int64 `json:"endedAtUnixTime"`
  ExitCode          int `json:"exitCode"`
}