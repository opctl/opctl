package dockercompose

import (
  "os/exec"
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core"
)

type opRunResourceFlusher interface {
  flush(
  pathToOpDir string,
  logChannel chan *models.LogEntry,
  ) (err error)
}

func newOpRunResourceFlusher(
) opRunResourceFlusher {

  return &_opRunResourceFlusher{}

}

type _opRunResourceFlusher struct{}

func (this _opRunResourceFlusher) flush(
pathToOpDir string,
logChannel chan *models.LogEntry,
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

  dockerComposeDownCmd.Stdout = core.NewLogEmittingIoWriter(logChannel, models.StdOutStream)
  dockerComposeDownCmd.Stderr = core.NewLogEmittingIoWriter(logChannel, models.StdErrStream)

  err = dockerComposeDownCmd.Run()

  return

}
