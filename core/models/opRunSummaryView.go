package models

import "time"

func NewOpRunSummaryView(
id                string,
opUrl     *Url,
startTime time.Time,
finishTime   time.Time,
exitCode           int,
) *OpRunSummaryView {

  return &OpRunSummaryView{
    Id:id,
    OpUrl:opUrl,
    StartTime:startTime,
    FinishTime:finishTime,
    ExitCode:exitCode,
  }

}

type OpRunSummaryView struct {
  Id         string `json:"id"`
  OpUrl      *Url `json:"opUrl"`
  StartTime  time.Time `json:"startTime"`
  FinishTime time.Time `json:"finishTime"`
  ExitCode   int `json:"exitCode"`
}