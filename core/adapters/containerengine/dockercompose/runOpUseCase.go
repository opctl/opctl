package dockercompose

//go:generate counterfeiter -o ./fakeRunOpUseCase.go --fake-name fakeRunOpUseCase ./ runOpUseCase

import (
  "os/exec"
  "github.com/opctl/engine/core/logging"
  "fmt"
  "github.com/opctl/engine/core/models"
  "time"
  "errors"
)

type runOpUseCase interface {
  Execute(
  correlationId string,
  opArgs map[string]string,
  opBundlePath string,
  opName string,
  opRunId string,
  logger logging.Logger,
  ) (err error)
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
opArgs map[string]string,
opBundlePath string,
opName string,
opRunId string,
logger logging.Logger,
) (err error) {

  // up
  dockerComposeUpCmd :=
  exec.Command(
    "docker-compose",
    "-p",
    opRunId,
    "up",
    "--force-recreate",
    "--abort-on-container-exit",
    "--remove-orphans",
    "--build",
  )

  dockerComposeUpCmd.Dir = opBundlePath

  dockerComposeUpCmd.Stdout = logging.NewLoggableIoWriter(correlationId, logging.StdOutStream, logger)
  dockerComposeUpCmd.Stderr = logging.NewLoggableIoWriter(correlationId, logging.StdErrStream, logger)

  for argName, argVal := range opArgs {
    dockerComposeUpCmd.Env = append(
      dockerComposeUpCmd.Env,
      fmt.Sprintf("%v=%v", argName, argVal),
    )
  }

  runError := dockerComposeUpCmd.Run()
  if (nil != runError) {
    logger(
      models.NewLogEntryEmittedEvent(
        correlationId,
        time.Now().UTC(),
        runError.Error(),
        logging.StdErrStream,
      ),
    )
  }

  exitCode, exitCodeReadError := this.opRunExitCodeReader.read(
    opBundlePath,
    opName,
    opRunId,
  )

  if (nil != exitCodeReadError) {
    err = errors.New(fmt.Sprintf("unable to read op exit code. Error was: %v", exitCodeReadError))
  } else if ( 0 != exitCode) {
    err = errors.New(fmt.Sprintf("nonzero op exit code. Exit code was: %v", exitCode))
  }

  defer func() {
    flushError := this.opRunResourceFlusher.flush(
      correlationId,
      opArgs,
      opBundlePath,
      opRunId,
      logger,
    )
    if (nil != flushError) {
      logger(
        models.NewLogEntryEmittedEvent(
          correlationId,
          time.Now().UTC(),
          flushError.Error(),
          logging.StdErrStream,
        ),
      )
    }

  }()

  return

}
