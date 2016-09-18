package docker

//go:generate counterfeiter -o ./fakeContainerRemover.go --fake-name fakeContainerRemover ./ containerRemover

import (
  "os/exec"
)

type containerRemover interface {
  // ensureContainerRemoved ensures any container associated with the provided opBundlePath and opRunId are removed
  EnsureContainerRemoved(
  opBundlePath string,
  opRunId string,
  )
}

func newContainerRemover(
) containerRemover {

  return &_containerRemover{}

}

type _containerRemover struct{}

func (this _containerRemover) EnsureContainerRemoved(
opBundlePath string,
opRunId string,
) {

  dockerComposeDownCmd :=
    exec.Command(
      "docker-compose",
      "-p",
      opRunId,
      "down",
      "--rmi",
      "local",
      "-v",
      "--remove-orphans",
    )

  dockerComposeDownCmd.Dir = opBundlePath

  dockerComposeDownCmd.Run()

  return

}
