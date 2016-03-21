package dockercompose

import (
  "os/exec"
  "os"
  "strings"
  dockerEngine "github.com/docker/engine-api/client"
  "golang.org/x/net/context"
)

type devOpExitCodeReader interface {
  read(
  devOpName string,
  ) (devOpExitCode int, err error)
}

func newDevOpExitCodeReader(
fs filesystem,
dockerEngine *dockerEngine.Client,
) devOpExitCodeReader {

  return &devOpExitCodeReaderImpl{
    fs:fs,
    dockerEngine:dockerEngine,
  }

}

type devOpExitCodeReaderImpl struct {
  fs           filesystem
  dockerEngine *dockerEngine.Client
}

func (r devOpExitCodeReaderImpl) read(
devOpName string,
) (devOpExitCode int, err error) {

  var relPathToDevOpDockerComposeFile string
  relPathToDevOpDockerComposeFile, err = r.fs.getRelPathToDevOpDockerComposeFile(devOpName)
  if (nil != err) {
    return
  }

  // get container id
  dockerComposePsCmd :=
  exec.Command(
    "docker-compose",
    "-f",
    relPathToDevOpDockerComposeFile,
    "ps",
    "-q",
    devOpName,
  )
  dockerComposePsCmd.Stderr = os.Stderr

  var dockerComposePsCmdRawOutput []byte
  dockerComposePsCmdRawOutput, err = dockerComposePsCmd.Output()
  if (nil != err) {
    return
  }

  container, err := r.dockerEngine.ContainerInspect(context.Background(), strings.TrimSpace(string(dockerComposePsCmdRawOutput)))
  if (nil != err) {
    return
  }
  devOpExitCode = container.State.ExitCode

  return

}
