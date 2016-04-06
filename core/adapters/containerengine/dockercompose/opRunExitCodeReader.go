package dockercompose

import (
  "os/exec"
  "os"
  "strings"
  dockerEngine "github.com/docker/engine-api/client"
  "golang.org/x/net/context"
)

type opRunExitCodeReader interface {
  read(
  pathToOpDir string,
  opName string,
  ) (opExitCode int, err error)
}

func newOpRunExitCodeReader(
fs filesystem,
dockerEngine *dockerEngine.Client,
) opRunExitCodeReader {

  return &_opRunExitCodeReader{
    fs:fs,
    dockerEngine:dockerEngine,
  }

}

type _opRunExitCodeReader struct {
  fs           filesystem
  dockerEngine *dockerEngine.Client
}

func (this _opRunExitCodeReader) read(
pathToOpDir string,
opName string,
) (opExitCode int, err error) {

  // get container id
  dockerComposePsCmd :=
  exec.Command(
    "docker-compose",
    "ps",
    "-q",
    opName,
  )

  dockerComposePsCmd.Dir = pathToOpDir

  dockerComposePsCmd.Stderr = os.Stderr

  var dockerComposePsCmdRawOutput []byte
  dockerComposePsCmdRawOutput, err = dockerComposePsCmd.Output()
  if (nil != err) {
    return
  }

  container, err := this.dockerEngine.ContainerInspect(context.Background(), strings.TrimSpace(string(dockerComposePsCmdRawOutput)))
  if (nil != err) {
    return
  }
  opExitCode = container.State.ExitCode

  return

}
