package dockercompose

//go:generate counterfeiter -o ./fakeOpRunResourceFlusher.go --fake-name fakeOpRunResourceFlusher ./ opRunResourceFlusher

import (
  "os/exec"
  "github.com/opctl/engine/core/logging"
  "fmt"
)

type opRunResourceFlusher interface {
  flush(
  correlationId string,
  opArgs map[string]string,
  opBundlePath string,
  opNamespace string,
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
opArgs map[string]string,
opBundlePath string,
opNamespace string,
logger logging.Logger,
) (err error) {

  // down
  dockerComposeDownCmd :=
  exec.Command(
    "docker-compose",
    "-p",
    opNamespace,
    "down",
    "--rmi",
    "local",
    "-v",
    "--remove-orphans",
  )

  dockerComposeDownCmd.Dir = opBundlePath

  dockerComposeDownCmd.Stdout = logging.NewLoggableIoWriter(correlationId, logging.StdOutStream, logger)
  dockerComposeDownCmd.Stderr = logging.NewLoggableIoWriter(correlationId, logging.StdErrStream, logger)

  for argName, argVal := range opArgs {
    dockerComposeDownCmd.Env = append(
      dockerComposeDownCmd.Env,
      fmt.Sprintf("%v=%v", argName, argVal),
    )
  }

  err = dockerComposeDownCmd.Run()

  return

}
