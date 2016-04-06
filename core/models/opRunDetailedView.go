package models

type OpRunDetailedView struct {
  Id                *string `json:"id"`
  OpUrl      *Url `json:"opUrl"`
  SubOps     []*OpRunSummaryView `json:"subOps,omitempty"`
  StartedAtUnixTime int64 `json:"startedAtUnixTime"`
  EndedAtUnixTime   int64 `json:"endedAtUnixTime"`
  ExitCode          int `json:"exitCode"`
}