package dockercompose

import (
  "os"
  "os/exec"
)

type opRunResourceFlusher interface {
  flush(
  pathToOpDir string,
  ) (err error)
}

func newOpRunResourceFlusher(
fs filesystem,
) opRunResourceFlusher {

  return &_opRunResourceFlusher{}

}

type _opRunResourceFlusher struct{}

func (this _opRunResourceFlusher) flush(
pathToOpDir string,
) (err error) {

  // down
  dockerComposeDownCmd :=
  exec.Command(
    "docker-compose",
    "down",
    "--rmi",
    "local",
    "-v",
  )

  dockerComposeDownCmd.Dir = pathToOpDir

  dockerComposeDownCmd.Stdout = os.Stdout
  dockerComposeDownCmd.Stderr = os.Stderr

  err = dockerComposeDownCmd.Run()

  return

}
