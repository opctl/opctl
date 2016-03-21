package models

type PipelineRunView struct {
  pipelineName       string
  stageRuns          []PipelineStageRunView
  startedAtPosixTime int64
  endedAtPosixTime   int64
  exitCode           int
}

func newPipelineRunView(
pipelineName       string,
stageRuns          []PipelineStageRunView,
startedAtPosixTime int64,
endedAtPosixTime   int64,
exitCode           int,
) PipelineRunView {

  return PipelineRunView{
    pipelineName:pipelineName,
    stageRuns:stageRuns,
    startedAtPosixTime:startedAtPosixTime,
    endedAtPosixTime:endedAtPosixTime,
    exitCode:exitCode,
  }

}

func (pipelineRunView PipelineRunView) PipelineName() string {
  return pipelineRunView.pipelineName
}

func (pipelineRunView PipelineRunView) StageRuns() []PipelineStageRunView {
  return pipelineRunView.stageRuns
}

func (pipelineRunView PipelineRunView) StartedAtPosixTime() int64 {
  return pipelineRunView.startedAtPosixTime
}

func (pipelineRunView PipelineRunView) EndedAtPosixTime() int64 {
  return pipelineRunView.endedAtPosixTime
}

func (pipelineRunView PipelineRunView) ExitCode() int {
  return pipelineRunView.exitCode
}