package models

type pipelineRunView struct {
  PipelineName       string `json:"pipelineName"`
  StageRuns          []PipelineStageRunView `json:"stageRuns"`
  StartedAtEpochTime int64 `json:"startedAtEpochTime"`
  EndedAtEpochTime   int64 `json:"endedAtEpochTime"`
  ExitCode           int `json:"exitCode"`
}

type PipelineRunView struct {
  pipelineRunView
}

func newPipelineRunView(
pipelineName       string,
stageRuns          []PipelineStageRunView,
startedAtEpochTime int64,
endedAtEpochTime   int64,
exitCode           int,
) PipelineRunView {

  return PipelineRunView{
    pipelineRunView{
      PipelineName:pipelineName,
      StageRuns:stageRuns,
      StartedAtEpochTime:startedAtEpochTime,
      EndedAtEpochTime:endedAtEpochTime,
      ExitCode:exitCode,
    },
  }

}

func (this PipelineRunView) PipelineName() string {
  return this.pipelineRunView.PipelineName
}

func (this PipelineRunView) StageRuns() []PipelineStageRunView {
  return this.pipelineRunView.StageRuns
}

func (this PipelineRunView) StartedAtEpochTime() int64 {
  return this.pipelineRunView.StartedAtEpochTime
}

func (this PipelineRunView) EndedAtEpochTime() int64 {
  return this.pipelineRunView.EndedAtEpochTime
}

func (this PipelineRunView) ExitCode() int {
  return this.pipelineRunView.ExitCode
}