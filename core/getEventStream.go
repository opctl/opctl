package core

import "github.com/opspec-io/sdk-golang/pkg/models"

func (this _core) GetEventStream(
subscriberEventChannel chan models.Event,
) (err error) {

  this.eventStream.RegisterSubscriber(subscriberEventChannel)

  return
}
