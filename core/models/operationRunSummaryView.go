package models

func NewOperationRunSummaryView(
id                *string,
operationName     string,
startedAtUnixTime int64,
endedAtUnixTime   int64,
exitCode           int,
) *OperationRunSummaryView {

  return &OperationRunSummaryView{
    Id:id,
    OperationName:operationName,
    StartedAtUnixTime:startedAtUnixTime,
    EndedAtUnixTime:endedAtUnixTime,
    ExitCode:exitCode,
  }

}

type OperationRunSummaryView struct {
  Id                *string `json:"id"`
  OperationName     string `json:"operationName"`
  StartedAtUnixTime int64 `json:"startedAtUnixTime"`
  EndedAtUnixTime   int64 `json:"endedAtUnixTime"`
  ExitCode          int `json:"exitCode"`
}