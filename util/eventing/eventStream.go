package eventing

import (
  "github.com/opspec-io/sdk-golang/pkg/model"
  "time"
)

func NewEventStream(
) EventStream {

  objectUnderConstruction := &eventStream{
    pendingPublishesChannel:make(chan model.Event),
    pendingSubscribesChannel:make(chan chan model.Event),
    pendingUnsubscribesChannel:make(chan chan model.Event),
    subscribers:make(map[chan model.Event]bool),
  }

  go objectUnderConstruction.init()

  return objectUnderConstruction

}

type EventPublisher interface {
  Publish(
  event model.Event,
  )
}

type EventStream interface {
  EventPublisher

  RegisterSubscriber(
  eventChannel chan model.Event,
  )

  UnregisterSubscriber(
  eventChannel chan model.Event,
  )
}

type eventStream struct {
  pendingPublishesChannel    chan model.Event

  pendingSubscribesChannel   chan chan model.Event

  pendingUnsubscribesChannel chan chan model.Event

  subscribers                map[chan model.Event]bool
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
eventChannel chan model.Event,
) {

  this.pendingSubscribesChannel <- eventChannel

}

func (this *eventStream) Publish(
event model.Event,
) {

  this.pendingPublishesChannel <- event

}

func (this *eventStream)   UnregisterSubscriber(
eventChannel chan model.Event,
) {

  this.pendingUnsubscribesChannel <- eventChannel

}
