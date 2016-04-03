package dockercompose

import (
  "os"
  "os/signal"
  "os/exec"
  "errors"
  "fmt"
  "time"
  "syscall"
  "github.com/dev-op-spec/engine/core/models"
  "path/filepath"
)

type runOperationUseCase interface {
  Execute(
  pathToOperationDir string,
  ) (operationRun models.OperationRunView, err error)
}

func newRunOperationUseCase(
fs filesystem,
operationRunExitCodeReader operationRunExitCodeReader,
operationRunResourceFlusher operationRunResourceFlusher,
) runOperationUseCase {

  return &_runOperationUseCase{
    fs:fs,
    operationRunExitCodeReader: operationRunExitCodeReader,
    operationRunResourceFlusher: operationRunResourceFlusher,
  }

}

type _runOperationUseCase struct {
  fs                          filesystem
  operationRunExitCodeReader  operationRunExitCodeReader
  operationRunResourceFlusher operationRunResourceFlusher
}

func (this _runOperationUseCase) Execute(
pathToOperationDir string,
) (operationRunView models.OperationRunView, err error) {

  operationRunView.StartedAtUnixTime = time.Now().Unix()
  operationRunView.OperationName = filepath.Base(pathToOperationDir)

  pathToOperationDockerComposeFile := this.fs.getPathToOperationDockerComposeFile(pathToOperationDir)

  // up
  dockerComposeUpCmd :=
  exec.Command(
    "docker-compose",
    "-f",
    pathToOperationDockerComposeFile,
    "up",
    "--force-recreate",
    "--abort-on-container-exit",
  )
  dockerComposeUpCmd.Stdout = os.Stdout
  dockerComposeUpCmd.Stderr = os.Stderr
  dockerComposeUpCmd.Stdin = os.Stdin

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

    operationRunView.ExitCode, err = this.operationRunExitCodeReader.read(
      pathToOperationDockerComposeFile,
    )
    if (0 != operationRunView.ExitCode) {

      runError := errors.New(fmt.Sprintf("%v exit code was: %v", operationRunView.OperationName, operationRunView.ExitCode))
      if (nil == err) {
        err = runError
      }else {
        err = errors.New(err.Error() + "\n" + runError.Error())
      }

    }

    flushOperationRunResourcesError := this.operationRunResourceFlusher.flush(
      pathToOperationDockerComposeFile,
    )
    if (nil != flushOperationRunResourcesError) {

      if (nil == err) {
        err = flushOperationRunResourcesError
      } else {
        err = errors.New(err.Error() + "\n" + flushOperationRunResourcesError.Error())
      }

      operationRunView.ExitCode = 1

    }

    // send resourceFlushIsCompleteChannel
    resourceFlushIsCompleteChannel <- true

    operationRunView.EndedAtUnixTime = time.Now().Unix()

  }()

  err = dockerComposeUpCmd.Run()

  return

}