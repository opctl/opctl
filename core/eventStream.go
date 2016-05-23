package core

import (
  "github.com/opctl/engine/core/models"
  "sync"
)

func newEventStream(
) eventStream {

  return &_eventStream{}

}

type eventStream interface {
  RegisterSubscriber(
  eventChannel chan models.Event,
  )

  Publish(
  event models.Event,
  )
}

type _eventStream struct {
  cachedEventsRWMutex sync.RWMutex
  cachedEvents        []models.Event

  subscribersRWMutex  sync.RWMutex
  subscribers         []chan models.Event
}

func (this *_eventStream) RegisterSubscriber(
eventChannel chan models.Event,
) {

  // don't block; subscriber may not be ready
  go func() {

    this.subscribersRWMutex.Lock()
    if (nil == this.subscribers) {
      // handle first subscriber
      this.subscribers = []chan models.Event{eventChannel}
    } else {
      this.subscribers = append(this.subscribers, eventChannel)
    }
    this.subscribersRWMutex.Unlock()

    this.cachedEventsRWMutex.RLock()
    for _, event := range this.cachedEvents {
      // return cached
      eventChannel <- event
    }
    this.cachedEventsRWMutex.RUnlock()

  }()

}

func (this *_eventStream) Publish(
event models.Event,
) {

  this.cachedEventsRWMutex.Lock()
  this.cachedEvents = append(this.cachedEvents, event)
  this.cachedEventsRWMutex.Unlock()

  this.subscribersRWMutex.RLock()
  subscribers := this.subscribers
  this.subscribersRWMutex.RUnlock()

  for _, eventChannel := range subscribers {

    // don't block; subscriber may not be ready
    go func(eventChannel chan models.Event) {

      eventChannel <- event

    }(eventChannel)
  }

}
