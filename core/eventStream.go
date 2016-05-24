package core

import (
  "github.com/opctl/engine/core/models"
  "time"
)

func newEventStream(
) eventStream {

  objectUnderConstruction := &_eventStream{
    pendingPublishesChannel:make(chan models.Event),
    pendingSubscribesChannel:make(chan chan models.Event),
    pendingUnsubscribesChannel:make(chan chan models.Event),
    subscribers:make(map[chan models.Event]bool),
  }

  go objectUnderConstruction.init()

  return objectUnderConstruction

}

type eventStream interface {
  RegisterSubscriber(
  eventChannel chan models.Event,
  )

  Publish(
  event models.Event,
  )

  UnregisterSubscriber(
  eventChannel chan models.Event,
  )
}

type _eventStream struct {
  pendingPublishesChannel    chan models.Event

  pendingSubscribesChannel   chan chan models.Event

  pendingUnsubscribesChannel chan chan models.Event

  subscribers                map[chan models.Event]bool
}

func (this *_eventStream) init(
) {

  for {

    select {

    case event := <-this.pendingPublishesChannel:
      for subscriber := range this.subscribers {
        select {
        case subscriber <- event:
        case <-time.After(time.Second * 5):
          delete(this.subscribers, subscriber)
          close(subscriber)
        }
      }

    case subscribe := <-this.pendingSubscribesChannel:
      this.subscribers[subscribe] = true

    case unsubscribe := <-this.pendingUnsubscribesChannel:
      if _, ok := this.subscribers[unsubscribe]; ok {
        delete(this.subscribers, unsubscribe)
        close(unsubscribe)
      }

    }

  }

}

func (this *_eventStream) RegisterSubscriber(
eventChannel chan models.Event,
) {

  this.pendingSubscribesChannel <- eventChannel

}

func (this *_eventStream) Publish(
event models.Event,
) {

  this.pendingPublishesChannel <- event

}

func (this *_eventStream)   UnregisterSubscriber(
eventChannel chan models.Event,
) {

  this.pendingUnsubscribesChannel <- eventChannel

}
