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

type runDevOpUcExecuter interface {
  Execute(
  devOpName string,
  ) (devOpRun models.DevOpRunView, err error)
}

func newRunDevOpUcExecuter(
fs filesystem,
ecr devOpExitCodeReader,
rf devOpResourceFlusher,
) runDevOpUcExecuter {

  return &runDevOpUcExecuterImpl{
    fs:fs,
    ecr: ecr,
    rf: rf,
  }

}

type runDevOpUcExecuterImpl struct {
  fs  filesystem
  ecr devOpExitCodeReader
  rf  devOpResourceFlusher
}

func (uc runDevOpUcExecuterImpl) Execute(
devOpName string,
) (devOpRunView models.DevOpRunView, err error) {

  devOpRunViewBuilder := models.NewDevOpRunViewBuilder()
  devOpRunViewBuilder.SetStartedAtPosixTime(time.Now().Unix())
  devOpRunViewBuilder.SetDevOpName(devOpName)

  var relPathToDevOpDockerComposeFile string
  relPathToDevOpDockerComposeFile, err = uc.fs.getRelPathToDevOpDockerComposeFile(devOpName)
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

    devOpRunViewBuilder.SetExitCode(130)

    // guard against hanging prompt
    os.Stdin.WriteString("\x03")

  }()

  defer func() {

    var devOpExitCode int
    devOpExitCode, err = uc.ecr.read(devOpName)
    if (0 != devOpExitCode) {

      runError := errors.New(fmt.Sprintf("%v exit code was: %v", devOpName, devOpExitCode))
      if (nil == err) {
        err = runError
      }else {
        err = errors.New(err.Error() + "\n" + runError.Error())
      }

      devOpRunViewBuilder.SetExitCode(devOpExitCode)

    }

    flushDevOpResourcesError := uc.rf.flush(devOpName)
    if (nil != flushDevOpResourcesError) {

      if (nil == err) {
        err = flushDevOpResourcesError
      } else {
        err = errors.New(err.Error() + "\n" + flushDevOpResourcesError.Error())
      }

      devOpRunViewBuilder.SetExitCode(1)

    }

    devOpRunViewBuilder.SetEndedAtPosixTime(time.Now().Unix())

    devOpRunView = devOpRunViewBuilder.Build()

  }()

  err = dockerComposeUpCmd.Run()

  return

}