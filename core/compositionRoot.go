package core

//go:generate counterfeiter -o ./fakeCompositionRoot.go --fake-name fakeCompositionRoot ./ compositionRoot

import (
  "github.com/opspec-io/sdk-golang/pkg/bundle"
)

type compositionRoot interface {
  StartOpRunUseCase() startOpRunUseCase
  GetEventStreamUseCase() getEventStreamUseCase
  KillOpRunUseCase() killOpRunUseCase
}

func newCompositionRoot(
containerEngine ContainerEngine,
) (compositionRoot compositionRoot) {

  /* factories */
  uniqueStringFactory := newUniqueStringFactory()

  /* components */
  eventStream := newEventStream()

  opspecSdk :=  bundle.New()

  storage := newStorage()

  opRunner := newOpRunner(
    containerEngine,
    eventStream,
    eventStream,
    opspecSdk,
    storage,
    uniqueStringFactory,
  )

  /* use cases */
  startOpRunUseCase := newStartOpRunUseCase(
    opRunner,
    eventStream,
    newPathNormalizer(),
    uniqueStringFactory,
  )

  getEventStreamUseCase := newGetEventStreamUseCase(
    eventStream,
  )

  killOpRunUseCase := newKillOpRunUseCase(
    containerEngine,
    eventStream,
    eventStream,
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
