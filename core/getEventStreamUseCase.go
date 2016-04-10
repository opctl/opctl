package core

import "github.com/dev-op-spec/engine/core/models"

func newGetEventStreamUseCase(
eventStream eventStream,
) getEventStreamUseCase {

  return &_getEventStreamUseCase{
    eventStream: eventStream,
  }

}

type getEventStreamUseCase interface {
  Execute(
  subscriberEventChannel chan models.Event,
  ) (err error)
}

type _getEventStreamUseCase struct {
  eventStream eventStream
}

func (this _getEventStreamUseCase) Execute(
subscriberEventChannel chan models.Event,
) (err error) {

  this.eventStream.RegisterSubscriber(subscriberEventChannel)

  return
}
