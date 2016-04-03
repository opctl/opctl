package dockercompose

import (
  "os"
  "os/exec"
)

type operationRunResourceFlusher interface {
  flush(
  pathToOperationDockerComposeFile string,
  ) (err error)
}

func newOperationRunResourceFlusher(
fs filesystem,
) operationRunResourceFlusher {

  return &_operationRunResourceFlusher{
    fs:fs,
  }

}

type _operationRunResourceFlusher struct {
  fs filesystem
}

func (this _operationRunResourceFlusher) flush(
pathToOperationDockerComposeFile string,
) (err error) {

  // down
  dockerComposeDownCmd :=
  exec.Command(
    "docker-compose",
    "-f",
    pathToOperationDockerComposeFile,
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
