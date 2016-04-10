package models

import "time"

type OpRunDetailedView struct {
  Id         string `json:"id"`
  OpUrl      *Url `json:"opUrl"`
  SubOps     []*OpRunSummaryView `json:"subOps,omitempty"`
  StartTime  time.Time `json:"startTime"`
  FinishTime time.Time `json:"finishTime"`
  ExitCode   int `json:"exitCode"`
}