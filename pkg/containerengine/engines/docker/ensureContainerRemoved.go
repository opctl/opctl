package docker

import (
  "os/exec"
)

func (this _containerEngine) EnsureContainerRemoved(
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
