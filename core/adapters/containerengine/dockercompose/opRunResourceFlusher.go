package dockercompose

import (
  "os/exec"
  "github.com/dev-op-spec/engine/core/logging"
)

type opRunResourceFlusher interface {
  flush(
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
  )

  dockerComposeDownCmd.Dir = pathToOpDir

  dockerComposeDownCmd.Stdout = logging.NewLoggableIoWriter(logging.StdOutStream, logger)
  dockerComposeDownCmd.Stderr = logging.NewLoggableIoWriter(logging.StdErrStream, logger)

  err = dockerComposeDownCmd.Run()

  return

}
