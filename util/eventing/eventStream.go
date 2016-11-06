package eventing

import (
  "github.com/opspec-io/sdk-golang/pkg/models"
  "time"
)

func NewEventStream(
) EventStream {

  objectUnderConstruction := &eventStream{
    pendingPublishesChannel:make(chan models.Event),
    pendingSubscribesChannel:make(chan chan models.Event),
    pendingUnsubscribesChannel:make(chan chan models.Event),
    subscribers:make(map[chan models.Event]bool),
  }

  go objectUnderConstruction.init()

  return objectUnderConstruction

}

type EventPublisher interface {
  Publish(
  event models.Event,
  )
}

type EventStream interface {
  EventPublisher

  RegisterSubscriber(
  eventChannel chan models.Event,
  )

  UnregisterSubscriber(
  eventChannel chan models.Event,
  )
}

type eventStream struct {
  pendingPublishesChannel    chan models.Event

  pendingSubscribesChannel   chan chan models.Event

  pendingUnsubscribesChannel chan chan models.Event

  subscribers                map[chan models.Event]bool
}

func (this *eventStream) init(
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

func (this *eventStream) RegisterSubscriber(
eventChannel chan models.Event,
) {

  this.pendingSubscribesChannel <- eventChannel

}

func (this *eventStream) Publish(
event models.Event,
) {

  this.pendingPublishesChannel <- event

}

func (this *eventStream)   UnregisterSubscriber(
eventChannel chan models.Event,
) {

  this.pendingUnsubscribesChannel <- eventChannel

}
