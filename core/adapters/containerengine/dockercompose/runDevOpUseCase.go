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
)

type runDevOpUseCase interface {
  Execute(
  devOpName string,
  ) (devOpRun models.DevOpRunView, err error)
}

func newRunDevOpUseCase(
fs filesystem,
ecr devOpExitCodeReader,
rf devOpResourceFlusher,
) runDevOpUseCase {

  return &_runDevOpUseCase{
    fs:fs,
    ecr: ecr,
    rf: rf,
  }

}

type _runDevOpUseCase struct {
  fs  filesystem
  ecr devOpExitCodeReader
  rf  devOpResourceFlusher
}

func (this _runDevOpUseCase) Execute(
devOpName string,
) (devOpRunView models.DevOpRunView, err error) {

  devOpRunView.StartedAtEpochTime = time.Now().Unix()
  devOpRunView.DevOpName = devOpName

  var relPathToDevOpDockerComposeFile string
  relPathToDevOpDockerComposeFile, err = this.fs.getRelPathToDevOpDockerComposeFile(devOpName)
  if (nil != err) {
    return
  }

  // up
  dockerComposeUpCmd :=
  exec.Command(
    "docker-compose",
    "-f",
    relPathToDevOpDockerComposeFile,
    "up",
    "--abort-on-container-exit",
  )
  dockerComposeUpCmd.Stdout = os.Stdout
  dockerComposeUpCmd.Stderr = os.Stderr
  dockerComposeUpCmd.Stdin = os.Stdin

  // handle SIGINT
  signalChannel := make(chan os.Signal, 1)
  signal.Notify(
    signalChannel,
    syscall.SIGHUP,
    syscall.SIGINT,
    syscall.SIGTERM,
    syscall.SIGQUIT,
  )

  // @TODO: this currently leaks memory if signal not received..
  go func() {
    <-signalChannel

    dockerComposeUpCmd.Process.Kill()

    devOpRunView.ExitCode = 130

    // guard against hanging prompt
    os.Stdin.WriteString("\x03")

  }()

  defer func() {

    var devOpExitCode int
    devOpExitCode, err = this.ecr.read(devOpName)
    if (0 != devOpExitCode) {

      runError := errors.New(fmt.Sprintf("%v exit code was: %v", devOpName, devOpExitCode))
      if (nil == err) {
        err = runError
      }else {
        err = errors.New(err.Error() + "\n" + runError.Error())
      }

      devOpRunView.ExitCode = devOpExitCode

    }

    flushDevOpResourcesError := this.rf.flush(devOpName)
    if (nil != flushDevOpResourcesError) {

      if (nil == err) {
        err = flushDevOpResourcesError
      } else {
        err = errors.New(err.Error() + "\n" + flushDevOpResourcesError.Error())
      }

      devOpRunView.ExitCode = 1

    }

    devOpRunView.EndedAtEpochTime = time.Now().Unix()

  }()

  err = dockerComposeUpCmd.Run()

  return

}