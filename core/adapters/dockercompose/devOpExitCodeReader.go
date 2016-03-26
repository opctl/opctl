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

  return &_devOpExitCodeReader{
    fs:fs,
    dockerEngine:dockerEngine,
  }

}

type _devOpExitCodeReader struct {
  fs           filesystem
  dockerEngine *dockerEngine.Client
}

func (this _devOpExitCodeReader) read(
devOpName string,
) (devOpExitCode int, err error) {

  var relPathToDevOpDockerComposeFile string
  relPathToDevOpDockerComposeFile, err = this.fs.getRelPathToDevOpDockerComposeFile(devOpName)
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

  container, err := this.dockerEngine.ContainerInspect(context.Background(), strings.TrimSpace(string(dockerComposePsCmdRawOutput)))
  if (nil != err) {
    return
  }
  devOpExitCode = container.State.ExitCode

  return

}
