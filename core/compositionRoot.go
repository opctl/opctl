package core

//go:generate counterfeiter -o ./fakeCompositionRoot.go --fake-name fakeCompositionRoot ./ compositionRoot

import (
  "github.com/opspec-io/engine/core/ports"
  "github.com/opspec-io/sdk-golang/models"
  "github.com/opspec-io/sdk-golang"
)

type compositionRoot interface {
  StartOpRunUseCase() startOpRunUseCase
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

  eventPublisher := func(logEntryEmittedEvent models.Event) {
    eventStream.Publish(logEntryEmittedEvent)
  }

  opspecSdk := opspec.New()

  storage := newStorage()

  opRunner := newOpRunner(
    containerEngine,
    eventStream,
    eventPublisher,
    opspecSdk,
    storage,
    uniqueStringFactory,
  )

  /* use cases */
  startOpRunUseCase := newStartOpRunUseCase(
    opRunner,
    eventPublisher,
    newPathNormalizer(),
    uniqueStringFactory,
  )

  getEventStreamUseCase := newGetEventStreamUseCase(
    eventStream,
  )

  killOpRunUseCase := newKillOpRunUseCase(
    containerEngine,
    eventStream,
    eventPublisher,
    storage,
    uniqueStringFactory,
  )

  compositionRoot = &_compositionRoot{
    startOpRunUseCase: startOpRunUseCase,
    getEventStreamUseCase:getEventStreamUseCase,
    killOpRunUseCase:killOpRunUseCase,
  }

  return

}

type _compositionRoot struct {
  startOpRunUseCase     startOpRunUseCase
  getEventStreamUseCase getEventStreamUseCase
  killOpRunUseCase      killOpRunUseCase
}

func (this _compositionRoot) StartOpRunUseCase() startOpRunUseCase {
  return this.startOpRunUseCase
}

func (this _compositionRoot) GetEventStreamUseCase() getEventStreamUseCase {
  return this.getEventStreamUseCase
}

func (this _compositionRoot) KillOpRunUseCase() killOpRunUseCase {
  return this.killOpRunUseCase
}
