package dockercompose

import (
  "os/exec"
  "github.com/open-devops/engine/core/logging"
)

type opRunResourceFlusher interface {
  flush(
  correlationId string,
  pathToOpDir string,
  logger logging.Logger,
  ) (err error)
}

func newOpRunResourceFlusher(
) opRunResourceFlusher {

  return &_opRunResourceFlusher{}

}

type _opRunResourceFlusher struct{}

func (this _opRunResourceFlusher) flush(
correlationId string,
pathToOpDir string,
logger logging.Logger,
) (err error) {

  // down
  dockerComposeDownCmd :=
  exec.Command(
    "docker-compose",
    "down",
    "--rmi",
    "local",
    "-v",
    "--remove-orphans",
  )

  dockerComposeDownCmd.Dir = pathToOpDir

  dockerComposeDownCmd.Stdout = logging.NewLoggableIoWriter(correlationId, logging.StdOutStream, logger)
  dockerComposeDownCmd.Stderr = logging.NewLoggableIoWriter(correlationId, logging.StdErrStream, logger)

  err = dockerComposeDownCmd.Run()

  return

}
