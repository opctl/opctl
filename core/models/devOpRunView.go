package models

type DevOpRunView struct {
  DevOpName          string
  StartedAtEpochTime int64
  EndedAtEpochTime   int64
  ExitCode           int
}

func newDevOpRunView(
devOpName          string,
startedAtEpochTime int64,
endedAtEpochTime   int64,
exitCode           int,
) DevOpRunView {

  return DevOpRunView{
    DevOpName:devOpName,
    StartedAtEpochTime:startedAtEpochTime,
    EndedAtEpochTime:endedAtEpochTime,
    ExitCode:exitCode,
  }

}