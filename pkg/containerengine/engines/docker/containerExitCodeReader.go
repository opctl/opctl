package docker

import (
  "os/exec"
  "strings"
  dockerEngine "github.com/docker/docker/client"
  "golang.org/x/net/context"
  "fmt"
)

type containerExitCodeReader interface {
  Read(
  opBundlePath string,
  opName string,
  opNamespace string,
  ) (opExitCode int, err error)
}

func newContainerExitCodeReader(
dockerEngine *dockerEngine.Client,
) containerExitCodeReader {

  return &_containerExitCodeReader{
    dockerEngine:dockerEngine,
  }

}

type _containerExitCodeReader struct {
  dockerEngine *dockerEngine.Client
}

func (this _containerExitCodeReader) Read(
opBundlePath string,
opName string,
opRunId string,
) (
opExitCode int,
err error,
) {

  // get container id
  dockerComposePsCmd :=
    exec.Command(
      "docker-compose",
      "-p",
      opRunId,
      "ps",
      "-q",
      opName,
    )

  dockerComposePsCmd.Dir = opBundlePath

  var dockerComposePsCmdRawOutput []byte
  dockerComposePsCmdRawOutput, dockerComposePsCmdErr := dockerComposePsCmd.Output()
  if (nil != dockerComposePsCmdErr) {
    switch dockerComposeRmCmdErr := dockerComposePsCmdErr.(type){
    case *exec.ExitError:
      err = fmt.Errorf("docker-compose ps returned error:\n  %v", string(dockerComposeRmCmdErr.Stderr))
    default:
      err = dockerComposeRmCmdErr
    }
    return
  }

  container, err := this.dockerEngine.ContainerInspect(
    context.Background(),
    strings.TrimSpace(string(dockerComposePsCmdRawOutput)),
  )
  if (nil != err) {
    return
  }
  opExitCode = container.State.ExitCode

  return

}
