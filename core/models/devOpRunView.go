package models

type DevOpRunView struct {
  DevOpName          string
  StartedAtUnixTime int64
  EndedAtUnixTime   int64
  ExitCode           int
}

func newDevOpRunView(
devOpName          string,
startedAtUnixTime int64,
endedAtUnixTime   int64,
exitCode           int,
) DevOpRunView {

  return DevOpRunView{
    DevOpName:devOpName,
    StartedAtUnixTime:startedAtUnixTime,
    EndedAtUnixTime:endedAtUnixTime,
    ExitCode:exitCode,
  }

}