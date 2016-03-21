package models

type DevOpRunViewBuilder interface {
  Build() DevOpRunView
  SetDevOpName(value string) DevOpRunViewBuilder
  SetStartedAtPosixTime(value int64) DevOpRunViewBuilder
  SetEndedAtPosixTime(value int64) DevOpRunViewBuilder
  SetExitCode(value int) DevOpRunViewBuilder
}

func NewDevOpRunViewBuilder() DevOpRunViewBuilder {
  return &devOpRunViewBuilder{}
}

type devOpRunViewBuilder struct {
  devOpName          string
  startedAtPosixTime int64
  endedAtPosixTime   int64
  exitCode           int
}

func (b *devOpRunViewBuilder) Build() DevOpRunView {

  return newDevOpRunView(
    b.devOpName,
    b.startedAtPosixTime,
    b.endedAtPosixTime,
    b.exitCode,
  )

}

func (b *devOpRunViewBuilder) SetDevOpName(value string) DevOpRunViewBuilder {

  b.devOpName = value
  return b

}

func (b *devOpRunViewBuilder) SetStartedAtPosixTime(value int64) DevOpRunViewBuilder {

  b.startedAtPosixTime = value
  return b

}

func (b *devOpRunViewBuilder) SetEndedAtPosixTime(value int64) DevOpRunViewBuilder {

  b.endedAtPosixTime = value
  return b

}

func (b *devOpRunViewBuilder) SetExitCode(value int) DevOpRunViewBuilder {

  b.exitCode = value
  return b

}
