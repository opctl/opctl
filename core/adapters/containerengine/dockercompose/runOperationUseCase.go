package dockercompose

import (
  "os"
  "os/signal"
  "os/exec"
  "errors"
  "fmt"
  "syscall"
  "path/filepath"
)

type runOperationUseCase interface {
  Execute(
  pathToOperationDir string,
  ) (exitCode int, err error)
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
) (exitCode int, err error) {

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

    exitCode, err = this.operationRunExitCodeReader.read(
      pathToOperationDockerComposeFile,
    )
    if (0 != exitCode) {

      runError := errors.New(
        fmt.Sprintf(
          "%v exit code was: %v",
          filepath.Base(pathToOperationDir),
          exitCode),
      )
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

      exitCode = 1

    }

    // send resourceFlushIsCompleteChannel
    resourceFlushIsCompleteChannel <- true

  }()

  err = dockerComposeUpCmd.Run()

  return

}