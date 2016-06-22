package dockercompose

//go:generate counterfeiter -o ./fakeKillOpRunUseCase.go --fake-name fakeKillOpRunUseCase ./ killOpRunUseCase

import (
  "github.com/opctl/engine/core/logging"
)

type killOpRunUseCase interface {
  Execute(
  correlationId string,
  opBundlePath string,
  opNamespace string,
  logger logging.Logger,
  ) (err error)
}

func newKillOpRunUseCase(
opRunResourceFlusher opRunResourceFlusher,
filesys filesys,
) killOpRunUseCase {

  return &_killOpRunUseCase{
    opRunResourceFlusher: opRunResourceFlusher,
    filesys:filesys,
  }

}

type _killOpRunUseCase struct {
  opRunResourceFlusher opRunResourceFlusher
  filesys           filesys
}

func (this _killOpRunUseCase) Execute(
correlationId string,
opBundlePath string,
opNamespace string,
logger logging.Logger,
) (err error) {

  if (this.filesys.isDockerComposeFileExistent(opBundlePath)) {

    this.opRunResourceFlusher.flush(
      correlationId,
      opBundlePath,
      opNamespace,
      logger,
    )

  }

  return

}
