package core

import (
  "errors"
  "time"
"github.com/dev-op-spec/engine/core/models"
"github.com/dev-op-spec/engine/core/ports"
)

type runPipelineUseCase interface {
  Execute(
  pipelineName string,
  namesOfAlreadyRunPipelines[]string,
  ) (pipelineRun models.PipelineRunView, err error)
}

func newRunPipelineUseCase(
fs ports.Filesys,
yml yamlCodec,
runDevOpUseCase runDevOpUseCase,
) runPipelineUseCase {

  return &_runPipelineUseCase{
    fs:fs,
    yml:yml,
    runDevOpUseCase: runDevOpUseCase,
  }

}

type _runPipelineUseCase struct {
  fs                 ports.Filesys
  yml                yamlCodec
  runDevOpUseCase runDevOpUseCase
}

func (this _runPipelineUseCase) Execute(
pipelineName string,
namesOfAlreadyRunPipelines[]string,
) (pipelineRun models.PipelineRunView, err error) {

  pipelineRunViewBuilder := models.NewPipelineRunViewBuilder()
  pipelineRunViewBuilder.SetStartedAtEpochTime(time.Now().Unix())
  pipelineRunViewBuilder.SetPipelineName(pipelineName)

  var pipelineFileBytes []byte
  pipelineFileBytes, err = this.fs.ReadPipelineFile(pipelineName)
  if (nil != err) {
    return
  }

  pipelineFile := pipelineFile{}
  err = this.yml.fromYaml(
    pipelineFileBytes,
    &pipelineFile,
  )
  if (nil != err) {
    return
  }

  defer func() {

    pipelineRunViewBuilder.SetEndedAtEpochTime(time.Now().Unix())

    pipelineRun = pipelineRunViewBuilder.Build()

  }()

  // guard infinite loop
  for _, previouslyRunPipeline := range namesOfAlreadyRunPipelines {

    if (previouslyRunPipeline == pipelineName) {
      err = errors.New("Unable to run pipeline with name=`" + pipelineName +
      "`. Pipelines cannot call themselves recursively.")
      return
    }

  }
  namesOfAlreadyRunPipelines = append(namesOfAlreadyRunPipelines, pipelineName)

  for _, stage := range pipelineFile.Stages {

    switch stage.Type {

    case devOpStageType:
      {

        var stageRun models.DevOpRunView

        stageRun, err = this.runDevOpUseCase.Execute(
          stage.Name,
        )

        pipelineRunViewBuilder.AddStageRun(stageRun)

        if (stageRun.ExitCode() != 0 || nil != err) {

          // bubble exit code up
          pipelineRunViewBuilder.SetExitCode(stageRun.ExitCode())
          return

        }

        break

      }

    case pipelineStageType:
      {

        var stageRun models.PipelineRunView

        stageRun, err = this.Execute(
          stage.Name,
          namesOfAlreadyRunPipelines,
        )

        pipelineRunViewBuilder.AddStageRun(stageRun)

        if (stageRun.ExitCode() != 0 || nil != err) {

          // bubble exit code up
          pipelineRunViewBuilder.SetExitCode(stageRun.ExitCode())
          return

        }

        break

      }

    default:
      {

        err = errors.New("Unable to run pipeline with name=`" + pipelineName +
        "`. Expected stage equal to `" + pipelineStageType + "` or `" + devOpStageType +
        "` but found `" + stage.Type + "`")

        return

      }

    }

  }

  return

}
