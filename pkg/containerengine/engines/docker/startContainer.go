package docker

//go:generate counterfeiter -o ./fakeStartContainerUseCase.go --fake-name fakeStartContainerUseCase ./ startContainerUseCase

import (
  "os/exec"
  "fmt"
  "errors"
  "github.com/opspec-io/engine/util/eventing"
)

func (this _containerEngine) StartContainer(
opRunArgs map[string]string,
opBundlePath string,
opName string,
opRunId string,
eventPublisher eventing.EventPublisher,
rootOpRunId string,
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

  dockerComposeUpCmd.Stdout = NewStdOutWriter(eventPublisher, opRunId, rootOpRunId)
  dockerComposeUpCmd.Stderr = NewStdErrWriter(eventPublisher, opRunId, rootOpRunId)

  for argName, argVal := range opRunArgs {
    dockerComposeUpCmd.Env = append(
      dockerComposeUpCmd.Env,
      fmt.Sprintf("%v=%v", argName, argVal),
    )
  }

  err = dockerComposeUpCmd.Run()
  if (nil != err) {
    return
  }

  exitCode, exitCodeReadError := this.containerExitCodeReader.Read(
    opBundlePath,
    opName,
    opRunId,
  )

  if (nil != exitCodeReadError) {
    err = errors.New(fmt.Sprintf("unable to read container exit code. Error was: %v", exitCodeReadError))
  } else if ( 0 != exitCode) {
    err = errors.New(fmt.Sprintf("nonzero container exit code. Exit code was: %v", exitCode))
  }

  defer this.EnsureContainerRemoved(
    opBundlePath,
    opRunId,
  )

  return

}
