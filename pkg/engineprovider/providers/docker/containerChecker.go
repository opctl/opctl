package docker

import (
  "fmt"
  "os/exec"
)

type containerChecker interface {
  IsValidContainerRunning(
  image string,
  ) (isContainerRunning bool, err error)
}

func newContainerChecker(
) (containerChecker containerChecker) {

  containerChecker = &_containerChecker{}

  return

}

type _containerChecker struct{}

func (this _containerChecker) IsValidContainerRunning(
image string,
) (isContainerRunning bool, err error) {

  dockerPsCmd :=
    exec.Command(
      "docker",
      "ps",
      "-q",
      "-f",
      fmt.Sprintf("name=%v", containerName),
      "-f",
      fmt.Sprintf("ancestor=%v", image),
    )

  dockerPsCmdOutput, dockerPsCmdErr := dockerPsCmd.Output()

  if (nil != dockerPsCmdErr) {
    switch dockerRmCmdErr := dockerPsCmdErr.(type){
    case *exec.ExitError:
      err = fmt.Errorf("Docker returned error:\n  %v", string(dockerRmCmdErr.Stderr))
    default:
      err = dockerRmCmdErr
    }
    return
  }

  if (len(dockerPsCmdOutput) > 0 ) {
    isContainerRunning = true
  }

  return
}
