package models

type PipelineRunViewBuilder interface {
  Build() PipelineRunView
  SetPipelineName(value string) PipelineRunViewBuilder
  AddStageRun(value PipelineStageRunView) PipelineRunViewBuilder
  SetStartedAtEpochTime(value int64) PipelineRunViewBuilder
  SetEndedAtEpochTime(value int64) PipelineRunViewBuilder
  SetExitCode(value int) PipelineRunViewBuilder
}

func NewPipelineRunViewBuilder() PipelineRunViewBuilder {
  return &pipelineRunViewBuilder{}
}

type pipelineRunViewBuilder struct {
  pipelineName       string
  stageRuns          []PipelineStageRunView
  startedAtEpochTime int64
  endedAtEpochTime   int64
  exitCode           int
}

func (b *pipelineRunViewBuilder) Build() PipelineRunView {

  return newPipelineRunView(
    b.pipelineName,
    b.stageRuns,
    b.startedAtEpochTime,
    b.endedAtEpochTime,
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

func (b *pipelineRunViewBuilder) SetStartedAtEpochTime(value int64) PipelineRunViewBuilder {

  b.startedAtEpochTime = value
  return b

}

func (b *pipelineRunViewBuilder) SetEndedAtEpochTime(value int64) PipelineRunViewBuilder {

  b.endedAtEpochTime = value
  return b

}

func (b *pipelineRunViewBuilder) SetExitCode(value int) PipelineRunViewBuilder {

  b.exitCode = value
  return b

}
