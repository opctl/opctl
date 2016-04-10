package dockercompose

import (
  "os"
  "os/signal"
  "os/exec"
  "errors"
  "fmt"
  "syscall"
  "github.com/dev-op-spec/engine/core/logging"
)

type runOpUseCase interface {
  Execute(
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
  )

  dockerComposeUpCmd.Dir = pathToOpDir

  // handle SIGINT
  signalChannel := make(chan os.Signal, 1)
  signal.Notify(
    signalChannel,
    syscall.SIGINT,
  )

  resourceFlushIsCompleteChannel := make(chan bool, 1)

  go func() {

    <-signalChannel

    // kill docker-compose; we flush our own resources
    dockerComposeUpCmd.Process.Kill()

    signal.Stop(signalChannel)

    // wait for resource flush to complete
    <-resourceFlushIsCompleteChannel

    if (nil != err) {

      fmt.Fprint(os.Stderr, err)

    }

    // exit with SIGINT exit code
    os.Exit(130)

  }()

  defer func() {

    exitCode, err = this.opRunExitCodeReader.read(
      pathToOpDir,
      opName,
    )
    if (0 != exitCode) {

      runError := errors.New(
        fmt.Sprintf(
          "%v exit code was: %v",
          opName,
          exitCode),
      )
      if (nil == err) {
        err = runError
      }else {
        err = errors.New(err.Error() + "\n" + runError.Error())
      }

    }

    flushOpRunResourcesError := this.opRunResourceFlusher.flush(
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

    // send resourceFlushIsComplete message
    resourceFlushIsCompleteChannel <- true

  }()

  dockerComposeUpCmd.Stdout = logging.NewLoggableIoWriter(logging.StdOutStream, logger)

  dockerComposeUpCmd.Stderr = logging.NewLoggableIoWriter(logging.StdErrStream, logger)

  err = dockerComposeUpCmd.Run()

  return

}