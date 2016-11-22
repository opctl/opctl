package core

import "github.com/opspec-io/sdk-golang/pkg/model"

func (this _core) GetEventStream(
subscriberEventChannel chan model.Event,
) (err error) {

  this.eventStream.RegisterSubscriber(subscriberEventChannel)

  return
}
