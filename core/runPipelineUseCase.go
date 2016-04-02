package core

import (
  "errors"
  "time"
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type runPipelineUseCase interface {
  Execute(
  req models.RunPipelineReq,
  namesOfAlreadyRunPipelines[]string,
  ) (pipelineRun models.PipelineRunView, err error)
}

func newRunPipelineUseCase(
filesys ports.Filesys,
pathToPipelineDirFactory pathToPipelineDirFactory,
pathToPipelineFileFactory pathToPipelineFileFactory,
runDevOpUseCase runDevOpUseCase,
uniqueStringFactory uniqueStringFactory,
yamlCodec yamlCodec,
) runPipelineUseCase {

  return &_runPipelineUseCase{
    filesys:filesys,
    pathToPipelineDirFactory:pathToPipelineDirFactory,
    pathToPipelineFileFactory:pathToPipelineFileFactory,
    runDevOpUseCase: runDevOpUseCase,
    uniqueStringFactory:uniqueStringFactory,
    yamlCodec:yamlCodec,
  }

}

type _runPipelineUseCase struct {
  filesys                   ports.Filesys
  pathToPipelineDirFactory  pathToPipelineDirFactory
  pathToPipelineFileFactory pathToPipelineFileFactory
  runDevOpUseCase           runDevOpUseCase
  uniqueStringFactory       uniqueStringFactory
  yamlCodec                 yamlCodec
}

func (this _runPipelineUseCase) Execute(
req models.RunPipelineReq,
namesOfAlreadyRunPipelines[]string,
) (pipelineRun models.PipelineRunView, err error) {

  pathToPipelineFile := this.pathToPipelineFileFactory.Construct(
    req.ProjectUrl,
    req.PipelineName,
  )

  pipelineRun.StartedAtUnixTime = time.Now().Unix()

  pipelineRun.Id, err = this.uniqueStringFactory.Construct()
  if (nil != err) {
    return
  }

  pipelineRun.PipelineName = req.PipelineName

  pipelineFileBytes, err := this.filesys.GetBytesOfFile(pathToPipelineFile)
  if (nil != err) {
    return
  }

  pipelineFile := pipelineFile{}
  err = this.yamlCodec.fromYaml(
    pipelineFileBytes,
    &pipelineFile,
  )
  if (nil != err) {
    return
  }

  defer func() {

    pipelineRun.EndedAtUnixTime = time.Now().Unix()

  }()

  // guard infinite loop
  for _, previouslyRunPipeline := range namesOfAlreadyRunPipelines {

    if (previouslyRunPipeline == req.PipelineName) {
      err = errors.New("Unable to run pipeline with name=`" + req.PipelineName +
      "`. Pipelines cannot call themselves recursively.")
      return
    }

  }
  namesOfAlreadyRunPipelines = append(namesOfAlreadyRunPipelines, req.PipelineName)

  for _, stage := range pipelineFile.Stages {

    switch stage.Type {

    case devOpStageType:
      {

        var devOpStageRun models.DevOpRunView
        devOpStageRun, err = this.runDevOpUseCase.Execute(
          *models.NewRunDevOpReq(
            req.ProjectUrl,
            stage.Name,
          ),
        )

        pipelineRun.Stages = append(
          pipelineRun.Stages,
          models.NewPipelineStageRunView(
            devOpStageRun.DevOpName,
            devOpStageType,
            devOpStageRun.StartedAtUnixTime,
            devOpStageRun.EndedAtUnixTime,
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
          *models.NewRunPipelineReq(
            req.ProjectUrl,
            stage.Name,
          ),
          namesOfAlreadyRunPipelines,
        )

        pipelineRun.Stages = append(
          pipelineRun.Stages,
          models.NewPipelineStageRunView(
            pipelineStageRun.PipelineName,
            pipelineStageType,
            pipelineStageRun.StartedAtUnixTime,
            pipelineStageRun.EndedAtUnixTime,
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

        err = errors.New("Unable to run pipeline with name=`" + req.PipelineName +
        "`. Expected stage equal to `" + pipelineStageType + "` or `" + devOpStageType +
        "` but found `" + stage.Type + "`")

        return

      }

    }

  }

  return

}
