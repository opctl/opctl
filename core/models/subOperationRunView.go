package models

func NewSubOperationRunView(
name               string,
startedAtUnixTime int64,
endedAtUnixTime   int64,
exitCode           int,
subOperations             []*SubOperationRunView,
) *SubOperationRunView {

  return &SubOperationRunView{
    Name:name,
    StartedAtUnixTime:startedAtUnixTime,
    EndedAtUnixTime:endedAtUnixTime,
    ExitCode:exitCode,
    SubOperations:subOperations,
  }

}

type SubOperationRunView struct {
  Name              string
  StartedAtUnixTime int64
  EndedAtUnixTime   int64
  ExitCode          int
  SubOperations     []*SubOperationRunView `json:",omitempty"`
}