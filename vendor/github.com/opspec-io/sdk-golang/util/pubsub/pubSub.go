package pubsub

//go:generate counterfeiter -o ./fakeEventSubscriber.go --fake-name FakeEventSubscriber ./ EventSubscriber
//go:generate counterfeiter -o ./fakeEventPublisher.go --fake-name FakeEventPublisher ./ EventPublisher
//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ PubSub

import (
	"github.com/opspec-io/sdk-golang/model"
	"sync"
	"time"
)

func New(
	eventRepo EventRepo,
) PubSub {
	return &pubSub{
		eventRepo:          eventRepo,
		subscriptions:      map[chan *model.Event]chan *model.Event{},
		subscriptionsMutex: sync.RWMutex{},
	}
}

type EventPublisher interface {
	Publish(
		event *model.Event,
	)
}

type EventSubscriber interface {
	Subscribe(
		filter *model.EventFilter,
		eventChannel chan *model.Event,
	)
}

type PubSub interface {
	EventPublisher
	EventSubscriber
}

type pubSub struct {
	eventRepo EventRepo
	// format: output : input
	subscriptions      map[chan *model.Event]chan *model.Event
	subscriptionsMutex sync.RWMutex
}

// O(n) complexity (n being topics subscribed to); thread safe
func (this *pubSub) Subscribe(
	filter *model.EventFilter,
	channel chan *model.Event,
) {
	this.subscriptionsMutex.Lock()
	defer this.subscriptionsMutex.Unlock()
	this.subscriptions[channel] = make(chan *model.Event, 10000)

	go func() {
		for _, event := range this.eventRepo.List(filter) {
			this.deliverEvent(event, channel)
		}
		for event := range this.subscriptions[channel] {
			RootOpId := getEventRootOpId(event)
			if !isRootOpIdExcludedByFilter(RootOpId, filter) {
				this.deliverEvent(event, channel)
			}
		}
	}()
}

func (this *pubSub) deliverEvent(
	event *model.Event,
	channel chan *model.Event,
) {
	select {
	case channel <- event:
	case <-time.After(time.Second * 5):
		this.subscriptionsMutex.Lock()
		delete(this.subscriptions, channel)
		this.subscriptionsMutex.Unlock()
		return
	}
}

// O(1) complexity; thread safe
func (this *pubSub) Publish(
	event *model.Event,
) {
	this.eventRepo.Add(event)

	this.subscriptionsMutex.RLock()
	defer this.subscriptionsMutex.RUnlock()
	for _, subscriptionInput := range this.subscriptions {
		subscriptionInput <- event
	}
}
