package core

import "github.com/opspec-io/sdk-golang/models"

func (this _api) GetEventStream(
subscriberEventChannel chan models.Event,
) (err error) {

  this.eventStream.RegisterSubscriber(subscriberEventChannel)

  return
}
