package dockercompose

import (
  "os/exec"
  "os"
  "strings"
  dockerEngine "github.com/docker/engine-api/client"
  "golang.org/x/net/context"
  "path"
)

type devOpRunExitCodeReader interface {
  read(
  pathToDevOpDockerComposeFile string,
  ) (devOpExitCode int, err error)
}

func newDevOpRunExitCodeReader(
fs filesystem,
dockerEngine *dockerEngine.Client,
) devOpRunExitCodeReader {

  return &_devOpRunExitCodeReader{
    fs:fs,
    dockerEngine:dockerEngine,
  }

}

type _devOpRunExitCodeReader struct {
  fs           filesystem
  dockerEngine *dockerEngine.Client
}

func (this _devOpRunExitCodeReader) read(
pathToDevOpDockerComposeFile string,
) (devOpExitCode int, err error) {

  devOpName := path.Base(path.Dir(pathToDevOpDockerComposeFile))

  // get container id
  dockerComposePsCmd :=
  exec.Command(
    "docker-compose",
    "-f",
    pathToDevOpDockerComposeFile,
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
