package docker

import (
  "fmt"
  "github.com/mitchellh/go-homedir"
  "os/exec"
)

type containerStarter interface {
  Start(
  image string,
  ) (err error)
}

func newContainerStarter(
pathNormalizer pathNormalizer,
) (containerStarter containerStarter) {

  containerStarter = _containerStarter{
    pathNormalizer:pathNormalizer,
  }

  return

}

type _containerStarter struct {
  pathNormalizer pathNormalizer
}

func (this _containerStarter) Start(
image string,
) (err error) {

  usersDir, err := homedir.Dir()
  if (nil != err) {
    return
  }

  normalizedUsersDir := this.pathNormalizer.Normalize(usersDir)

  dockerRunCmd :=
    exec.Command(
      "docker",
      "run",
      "-d",
      "-p",
      "42224:42224",
      "-v",
      fmt.Sprintf("%v:%v", normalizedUsersDir, normalizedUsersDir),
      "-v",
      "/var/run/docker.sock:/var/run/docker.sock",
      fmt.Sprintf("--name=%v", containerName),
      image,
    )

  _, dockerRunCmdErr := dockerRunCmd.Output()

  if (nil != dockerRunCmdErr) {
    switch dockerRunCmdErr := dockerRunCmdErr.(type){
    case *exec.ExitError:
      err = fmt.Errorf("Docker returned error:\n  %v", string(dockerRunCmdErr.Stderr))
    default:
      err = dockerRunCmdErr
    }
  }

  return
}
