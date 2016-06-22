package core

//go:generate counterfeiter -o ./fakeCompositionRoot.go --fake-name fakeCompositionRoot ./ compositionRoot

import (
  "github.com/opctl/engine/core/ports"
  "github.com/opctl/engine/core/models"
)

type compositionRoot interface {
  RunOpUseCase() runOpUseCase
  GetEventStreamUseCase() getEventStreamUseCase
  KillOpRunUseCase() killOpRunUseCase
}

func newCompositionRoot(
containerEngine ports.ContainerEngine,
filesys ports.Filesys,
) (compositionRoot compositionRoot, err error) {

  /* factories */
  uniqueStringFactory := newUniqueStringFactory()

  /* components */
  eventStream := newEventStream()

  yamlCodec := newYamlCodec()

  logger := func(logEntryEmittedEvent models.LogEntryEmittedEvent) {
    eventStream.Publish(logEntryEmittedEvent)
  }

  opRunner := newOpRunner(
    containerEngine,
    eventStream,
    filesys,
    logger,
    uniqueStringFactory,
    yamlCodec,
  )

  /* use cases */
  runOpUseCase := newRunOpUseCase(
    opRunner,
    uniqueStringFactory,
  )

  getEventStreamUseCase := newGetEventStreamUseCase(
    eventStream,
  )

  killOpRunUseCase := newKillOpRunUseCase(
    opRunner,
    uniqueStringFactory,
  )

  compositionRoot = &_compositionRoot{
    runOpUseCase: runOpUseCase,
    getEventStreamUseCase:getEventStreamUseCase,
    killOpRunUseCase:killOpRunUseCase,
  }

  return

}

type _compositionRoot struct {
  runOpUseCase          runOpUseCase
  getEventStreamUseCase getEventStreamUseCase
  killOpRunUseCase      killOpRunUseCase
}

func (this _compositionRoot) RunOpUseCase() runOpUseCase {
  return this.runOpUseCase
}

func (this _compositionRoot) GetEventStreamUseCase() getEventStreamUseCase {
  return this.getEventStreamUseCase
}

func (this _compositionRoot) KillOpRunUseCase() killOpRunUseCase {
  return this.killOpRunUseCase
}
