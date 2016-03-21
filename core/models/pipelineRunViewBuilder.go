package models

type PipelineRunViewBuilder interface {
  Build() PipelineRunView
  SetPipelineName(value string) PipelineRunViewBuilder
  AddStageRun(value PipelineStageRunView) PipelineRunViewBuilder
  SetStartedAtPosixTime(value int64) PipelineRunViewBuilder
  SetEndedAtPosixTime(value int64) PipelineRunViewBuilder
  SetExitCode(value int) PipelineRunViewBuilder
}

func NewPipelineRunViewBuilder() PipelineRunViewBuilder {
  return &pipelineRunViewBuilder{}
}

type pipelineRunViewBuilder struct {
  pipelineName       string
  stageRuns          []PipelineStageRunView
  startedAtPosixTime int64
  endedAtPosixTime   int64
  exitCode           int
}

func (b *pipelineRunViewBuilder) Build() PipelineRunView {

  return newPipelineRunView(
    b.pipelineName,
    b.stageRuns,
    b.startedAtPosixTime,
    b.endedAtPosixTime,
    b.exitCode,
  )

}

func (b *pipelineRunViewBuilder) SetPipelineName(value string) PipelineRunViewBuilder {

  b.pipelineName = value
  return b

}

func (b *pipelineRunViewBuilder) AddStageRun(value PipelineStageRunView) PipelineRunViewBuilder {

  b.stageRuns = append(b.stageRuns, value)
  return b

}

func (b *pipelineRunViewBuilder) SetStartedAtPosixTime(value int64) PipelineRunViewBuilder {

  b.startedAtPosixTime = value
  return b

}

func (b *pipelineRunViewBuilder) SetEndedAtPosixTime(value int64) PipelineRunViewBuilder {

  b.endedAtPosixTime = value
  return b

}

func (b *pipelineRunViewBuilder) SetExitCode(value int) PipelineRunViewBuilder {

  b.exitCode = value
  return b

}
