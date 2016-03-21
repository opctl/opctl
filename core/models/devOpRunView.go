package models

type DevOpRunView struct {
  devOpName          string
  startedAtPosixTime int64
  endedAtPosixTime   int64
  exitCode           int
}

func newDevOpRunView(
devOpName          string,
startedAtPosixTime int64,
endedAtPosixTime   int64,
exitCode           int,
) DevOpRunView {

  return DevOpRunView{
    devOpName:devOpName,
    startedAtPosixTime:startedAtPosixTime,
    endedAtPosixTime:endedAtPosixTime,
    exitCode:exitCode,
  }

}

func (devOpRunView DevOpRunView) DevOpName() string {
  return devOpRunView.devOpName
}

func (devOpRunView DevOpRunView) StartedAtPosixTime() int64 {
  return devOpRunView.startedAtPosixTime
}

func (devOpRunView DevOpRunView) EndedAtPosixTime() int64 {
  return devOpRunView.endedAtPosixTime
}

func (devOpRunView DevOpRunView) ExitCode() int {
  return devOpRunView.exitCode
}