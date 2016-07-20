package core

//go:generate counterfeiter -o ./fakeCompositionRoot.go --fake-name fakeCompositionRoot ./ compositionRoot

import (
  "github.com/opctl/engine/core/models"
  "github.com/opspec-io/sdk-golang"
  "github.com/opctl/engine/core/ports"
)

type compositionRoot interface {
  RunOpUseCase() runOpUseCase
  GetEventStreamUseCase() getEventStreamUseCase
  KillOpRunUseCase() killOpRunUseCase
}

func newCompositionRoot(
containerEngine ports.ContainerEngine,
) (compositionRoot compositionRoot) {

  /* factories */
  uniqueStringFactory := newUniqueStringFactory()

  /* components */
  eventStream := newEventStream()

  logger := func(logEntryEmittedEvent models.LogEntryEmittedEvent) {
    eventStream.Publish(logEntryEmittedEvent)
  }

  opspecSdk := opspec.New()

  storage := newStorage()

  opRunner := newOpRunner(
    containerEngine,
    eventStream,
    logger,
    opspecSdk,
    storage,
    uniqueStringFactory,
  )

  /* use cases */
  runOpUseCase := newRunOpUseCase(
    opRunner,
    logger,
    newPathNormalizer(),
    uniqueStringFactory,
  )

  getEventStreamUseCase := newGetEventStreamUseCase(
    eventStream,
  )

  killOpRunUseCase := newKillOpRunUseCase(
    containerEngine,
    eventStream,
    logger,
    storage,
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
