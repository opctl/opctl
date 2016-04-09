package models

func NewOpRunSummaryView(
id                string,
opUrl     *Url,
startedAtUnixTime int64,
endedAtUnixTime   int64,
exitCode           int,
) *OpRunSummaryView {

  return &OpRunSummaryView{
    Id:id,
    OpUrl:opUrl,
    StartedAtUnixTime:startedAtUnixTime,
    EndedAtUnixTime:endedAtUnixTime,
    ExitCode:exitCode,
  }

}

type OpRunSummaryView struct {
  Id                string `json:"id"`
  OpUrl      *Url `json:"opUrl"`
  StartedAtUnixTime int64 `json:"startedAtUnixTime"`
  EndedAtUnixTime   int64 `json:"endedAtUnixTime"`
  ExitCode          int `json:"exitCode"`
}