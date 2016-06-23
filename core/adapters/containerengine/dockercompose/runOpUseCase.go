package dockercompose

//go:generate counterfeiter -o ./fakeRunOpUseCase.go --fake-name fakeRunOpUseCase ./ runOpUseCase

import (
  "os/exec"
  "errors"
  "github.com/opctl/engine/core/logging"
  "fmt"
)

type runOpUseCase interface {
  Execute(
  correlationId string,
  opArgs map[string]string,
  opBundlePath string,
  opName string,
  opNamespace string,
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
opArgs map[string]string,
opBundlePath string,
opName string,
opNamespace string,
logger logging.Logger,
) (exitCode int, err error) {

  // up
  dockerComposeUpCmd :=
  exec.Command(
    "docker-compose",
    "-p",
    opNamespace,
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

  err = dockerComposeUpCmd.Run()
  if (nil != err) {
    exitCode = 1
  }

  exitCode, err = this.opRunExitCodeReader.read(
    opBundlePath,
    opName,
    opNamespace,
  )

  defer func() {

    flushOpRunResourcesError := this.opRunResourceFlusher.flush(
      correlationId,
      opArgs,
      opBundlePath,
      opNamespace,
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

  return

}
