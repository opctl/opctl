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

  this.cachedEventsRWMutex.RLock()
  for _, event := range this.cachedEvents {
    // return cached
    eventChannel <- event
  }
  this.cachedEventsRWMutex.RUnlock()

  this.subscribersRWMutex.Lock()
  if (nil == this.subscribers) {
    // handle first subscriber
    this.subscribers = []chan models.Event{eventChannel}
  }else {
    this.subscribers = append(this.subscribers, eventChannel)
  }
  this.subscribersRWMutex.Unlock()

}

func (this *_eventStream) Publish(
event models.Event,
) {

  this.cachedEventsRWMutex.Lock()
  this.cachedEvents = append(this.cachedEvents, event)
  this.cachedEventsRWMutex.Unlock()

  this.subscribersRWMutex.RLock()
  for _, eventChannel := range this.subscribers {
    eventChannel <- event
  }
  this.subscribersRWMutex.RUnlock()

  return

}
