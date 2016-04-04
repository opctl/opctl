package models

func NewOperationRunSummaryView(
id                *string,
operationUrl     *Url,
startedAtUnixTime int64,
endedAtUnixTime   int64,
exitCode           int,
) *OperationRunSummaryView {

  return &OperationRunSummaryView{
    Id:id,
    OperationUrl:operationUrl,
    StartedAtUnixTime:startedAtUnixTime,
    EndedAtUnixTime:endedAtUnixTime,
    ExitCode:exitCode,
  }

}

type OperationRunSummaryView struct {
  Id                *string `json:"id"`
  OperationUrl      *Url `json:"operationUrl"`
  StartedAtUnixTime int64 `json:"startedAtUnixTime"`
  EndedAtUnixTime   int64 `json:"endedAtUnixTime"`
  ExitCode          int `json:"exitCode"`
}