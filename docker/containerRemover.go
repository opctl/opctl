package docker

import (
  "os/exec"
  "fmt"
)

type containerRemover interface {
  RemoveIfExists(
  containerName string,
  ) (err error)
}

func newContainerRemover(
) (containerRemover containerRemover) {

  containerRemover = &_containerRemover{}

  return

}

type _containerRemover struct{}

func (this _containerRemover) RemoveIfExists(
containerName string,
) (err error) {

  dockerPsCmd :=
    exec.Command(
      "docker",
      "ps",
      "-a",
      "-q",
      "-f",
      fmt.Sprintf("name=%v", containerName),
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
    dockerRmCmd :=
      exec.Command(
        "docker",
        "rm",
        "-fv",
        containerName,
      )

    _, dockerRmCmdErr := dockerRmCmd.Output()

    if (nil != dockerRmCmdErr) {
      switch dockerRmCmdErr := dockerRmCmdErr.(type){
      case *exec.ExitError:
        err = fmt.Errorf("Docker returned error:\n  %v", string(dockerRmCmdErr.Stderr))
      default:
        err = dockerRmCmdErr
      }
    }
  }

  return
}
