package models

type DevOpRunViewBuilder interface {
  Build() DevOpRunView
  SetDevOpName(value string) DevOpRunViewBuilder
  SetStartedAtEpochTime(value int64) DevOpRunViewBuilder
  SetEndedAtEpochTime(value int64) DevOpRunViewBuilder
  SetExitCode(value int) DevOpRunViewBuilder
}

func NewDevOpRunViewBuilder() DevOpRunViewBuilder {
  return &devOpRunViewBuilder{}
}

type devOpRunViewBuilder struct {
  devOpName          string
  startedAtEpochTime int64
  endedAtEpochTime   int64
  exitCode           int
}

func (b *devOpRunViewBuilder) Build() DevOpRunView {

  return newDevOpRunView(
    b.devOpName,
    b.startedAtEpochTime,
    b.endedAtEpochTime,
    b.exitCode,
  )

}

func (b *devOpRunViewBuilder) SetDevOpName(value string) DevOpRunViewBuilder {

  b.devOpName = value
  return b

}

func (b *devOpRunViewBuilder) SetStartedAtEpochTime(value int64) DevOpRunViewBuilder {

  b.startedAtEpochTime = value
  return b

}

func (b *devOpRunViewBuilder) SetEndedAtEpochTime(value int64) DevOpRunViewBuilder {

  b.endedAtEpochTime = value
  return b

}

func (b *devOpRunViewBuilder) SetExitCode(value int) DevOpRunViewBuilder {

  b.exitCode = value
  return b

}
