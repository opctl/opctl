package dockercompose

import (
  "os/exec"
  "os"
  "strings"
  dockerEngine "github.com/docker/engine-api/client"
  "golang.org/x/net/context"
  "path"
)

type operationRunExitCodeReader interface {
  read(
  pathToOperationDockerComposeFile string,
  ) (operationExitCode int, err error)
}

func newOperationRunExitCodeReader(
fs filesystem,
dockerEngine *dockerEngine.Client,
) operationRunExitCodeReader {

  return &_operationRunExitCodeReader{
    fs:fs,
    dockerEngine:dockerEngine,
  }

}

type _operationRunExitCodeReader struct {
  fs           filesystem
  dockerEngine *dockerEngine.Client
}

func (this _operationRunExitCodeReader) read(
pathToOperationDockerComposeFile string,
) (operationExitCode int, err error) {

  operationName := path.Base(path.Dir(pathToOperationDockerComposeFile))

  // get container id
  dockerComposePsCmd :=
  exec.Command(
    "docker-compose",
    "-f",
    pathToOperationDockerComposeFile,
    "ps",
    "-q",
    operationName,
  )
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
  operationExitCode = container.State.ExitCode

  return

}
