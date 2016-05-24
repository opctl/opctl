package dockercompose

//go:generate counterfeiter -o ./fakeRunOpUseCase.go --fake-name fakeRunOpUseCase ./ runOpUseCase

import (
  "os/exec"
  "errors"
  "github.com/opctl/engine/core/logging"
)

type runOpUseCase interface {
  Execute(
  correlationId string,
  pathToOpDir string,
  opName string,
  logger logging.Logger,
  ) (exitCode int, err error)
}

func newRunOpUseCase(
opRunExitCodeReader opRunExitCodeReader,
opRunResourceFlusher opRunResourceFlusher,
) runOpUseCase {

  return &_runOpUseCase{
    opRunExitCodeReader: opRunExitCodeReader,
    opRunResourceFlusher: opRunResourceFlusher,
  }

}

type _runOpUseCase struct {
  opRunExitCodeReader  opRunExitCodeReader
  opRunResourceFlusher opRunResourceFlusher
}

func (this _runOpUseCase) Execute(
correlationId string,
pathToOpDir string,
opName string,
logger logging.Logger,
) (exitCode int, err error) {

  // up
  dockerComposeUpCmd :=
  exec.Command(
    "docker-compose",
    "up",
    "--force-recreate",
    "--abort-on-container-exit",
    "--remove-orphans",
    "--build",
  )

  dockerComposeUpCmd.Dir = pathToOpDir

  defer func() {

    flushOpRunResourcesError := this.opRunResourceFlusher.flush(
      correlationId,
      pathToOpDir,
      logger,
    )
    if (nil != flushOpRunResourcesError) {

      if (nil == err) {
        err = flushOpRunResourcesError
      } else {
        err = errors.New(err.Error() + "\n" + flushOpRunResourcesError.Error())
      }

      exitCode = 1

    }

  }()

  dockerComposeUpCmd.Stdout = logging.NewLoggableIoWriter(correlationId, logging.StdOutStream, logger)
  dockerComposeUpCmd.Stderr = logging.NewLoggableIoWriter(correlationId, logging.StdErrStream, logger)

  err = dockerComposeUpCmd.Run()
  if (nil != err) {
    exitCode = 1
  }

  exitCode, err = this.opRunExitCodeReader.read(
    opName,
    pathToOpDir,
  )

  return

}
