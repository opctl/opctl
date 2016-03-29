package dockercompose

import (
  "os"
  "os/exec"
)

type devOpRunResourceFlusher interface {
  flush(
  pathToDevOpDockerComposeFile string,
  ) (err error)
}

func newDevOpRunResourceFlusher(
fs filesystem,
) devOpRunResourceFlusher {

  return &_devOpRunResourceFlusher{
    fs:fs,
  }

}

type _devOpRunResourceFlusher struct {
  fs filesystem
}

func (this _devOpRunResourceFlusher) flush(
pathToDevOpDockerComposeFile string,
) (err error) {

  // down
  dockerComposeDownCmd :=
  exec.Command(
    "docker-compose",
    "-f",
    pathToDevOpDockerComposeFile,
    "down",
    "--rmi",
    "local",
    "-v",
  )
  dockerComposeDownCmd.Stdout = os.Stdout
  dockerComposeDownCmd.Stderr = os.Stderr
  err = dockerComposeDownCmd.Run()

  return

}
