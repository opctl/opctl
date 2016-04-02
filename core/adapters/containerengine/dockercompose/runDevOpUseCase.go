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

type runDevOpUseCase interface {
  Execute(
  pathToDevOpDir string,
  ) (devOpRun models.DevOpRunView, err error)
}

func newRunDevOpUseCase(
fs filesystem,
devOpRunExitCodeReader devOpRunExitCodeReader,
devOpRunResourceFlusher devOpRunResourceFlusher,
) runDevOpUseCase {

  return &_runDevOpUseCase{
    fs:fs,
    devOpRunExitCodeReader: devOpRunExitCodeReader,
    devOpRunResourceFlusher: devOpRunResourceFlusher,
  }

}

type _runDevOpUseCase struct {
  fs                      filesystem
  devOpRunExitCodeReader  devOpRunExitCodeReader
  devOpRunResourceFlusher devOpRunResourceFlusher
}

func (this _runDevOpUseCase) Execute(
pathToDevOpDir string,
) (devOpRunView models.DevOpRunView, err error) {

  devOpRunView.StartedAtUnixTime = time.Now().Unix()
  devOpRunView.DevOpName = filepath.Base(pathToDevOpDir)

  pathToDevOpDockerComposeFile := this.fs.getPathToDevOpDockerComposeFile(pathToDevOpDir)

  // up
  dockerComposeUpCmd :=
  exec.Command(
    "docker-compose",
    "-f",
    pathToDevOpDockerComposeFile,
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

    // exit with proper exit code
    os.Exit(130)

  }()

  defer func() {

    devOpRunView.ExitCode, err = this.devOpRunExitCodeReader.read(
      pathToDevOpDockerComposeFile,
    )
    if (0 != devOpRunView.ExitCode) {

      runError := errors.New(fmt.Sprintf("%v exit code was: %v", devOpRunView.DevOpName, devOpRunView.ExitCode))
      if (nil == err) {
        err = runError
      }else {
        err = errors.New(err.Error() + "\n" + runError.Error())
      }

    }

    flushDevOpRunResourcesError := this.devOpRunResourceFlusher.flush(
      pathToDevOpDockerComposeFile,
    )
    if (nil != flushDevOpRunResourcesError) {

      if (nil == err) {
        err = flushDevOpRunResourcesError
      } else {
        err = errors.New(err.Error() + "\n" + flushDevOpRunResourcesError.Error())
      }

      devOpRunView.ExitCode = 1

    }

    // send resourceFlushIsCompleteChannel
    resourceFlushIsCompleteChannel <- true

    devOpRunView.EndedAtUnixTime = time.Now().Unix()

  }()

  err = dockerComposeUpCmd.Run()

  return

}