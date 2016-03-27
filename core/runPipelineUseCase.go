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
  fs              ports.Filesys
  yml             yamlCodec
  runDevOpUseCase runDevOpUseCase
}

func (this _runPipelineUseCase) Execute(
pipelineName string,
namesOfAlreadyRunPipelines[]string,
) (pipelineRun models.PipelineRunView, err error) {

  pipelineRun.StartedAtEpochTime = time.Now().Unix()
  pipelineRun.PipelineName = pipelineName

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

    pipelineRun.EndedAtEpochTime = time.Now().Unix()

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

        var devOpStageRun models.DevOpRunView

        devOpStageRun, err = this.runDevOpUseCase.Execute(
          stage.Name,
        )

        pipelineRun.Stages = append(
          pipelineRun.Stages,
          models.NewPipelineStageRunView(
            devOpStageRun.DevOpName,
            devOpStageType,
            devOpStageRun.StartedAtEpochTime,
            devOpStageRun.EndedAtEpochTime,
            devOpStageRun.ExitCode,
            nil,
          ),
        )

        if (devOpStageRun.ExitCode != 0 || nil != err) {

          // bubble exit code up
          pipelineRun.ExitCode = devOpStageRun.ExitCode
          return

        }

        break

      }

    case pipelineStageType:
      {

        var pipelineStageRun models.PipelineRunView

        pipelineStageRun, err = this.Execute(
          stage.Name,
          namesOfAlreadyRunPipelines,
        )

        pipelineRun.Stages = append(
          pipelineRun.Stages,
          models.NewPipelineStageRunView(
            pipelineStageRun.PipelineName,
            pipelineStageType,
            pipelineStageRun.StartedAtEpochTime,
            pipelineStageRun.EndedAtEpochTime,
            pipelineStageRun.ExitCode,
            nil,
          ),
        )

        if (pipelineStageRun.ExitCode != 0 || nil != err) {

          // bubble exit code up
          pipelineRun.ExitCode = pipelineStageRun.ExitCode
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
