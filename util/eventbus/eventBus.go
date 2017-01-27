package eventbus

//go:generate counterfeiter -o ./fakeEventPublisher.go --fake-name FakeEventPublisher ./ EventPublisher
//go:generate counterfeiter -o ./fakeEventBus.go --fake-name FakeEventBus ./ EventBus

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"sync"
	"time"
)

func New() EventBus {

	objectUnderConstruction := &eventBus{
		pendingPublishesChannel: make(chan model.Event),
		subscriberMutex:         sync.RWMutex{},
		subscribers:             make(map[chan model.Event]*model.EventFilter),
	}

	go objectUnderConstruction.init()

	return objectUnderConstruction

}

// publishes messages to the event bus
type EventPublisher interface {
	Publish(
		event model.Event,
	)
}

type EventBus interface {
	EventPublisher

	RegisterSubscriber(
		filter *model.EventFilter,
		eventChannel chan model.Event,
	)

	UnregisterSubscriber(
		eventChannel chan model.Event,
	)
}

type eventBus struct {
	pendingPublishesChannel chan model.Event
	subscriberMutex         sync.RWMutex
	subscribers             map[chan model.Event]*model.EventFilter
}

func (this *eventBus) init() {

	for {
		event := <-this.pendingPublishesChannel
		this.publishToSubscribers(event)
	}

}

// O(1) complexity; thread safe
func (this *eventBus) RegisterSubscriber(
	filter *model.EventFilter,
	eventChannel chan model.Event,
) {
	this.subscriberMutex.Lock()
	defer this.subscriberMutex.Unlock()
	this.subscribers[eventChannel] = filter
}

func (this *eventBus) isEventFiltered(
	event model.Event,
	filter *model.EventFilter,
) bool {
	if nil != filter && nil != filter.OpGraphIds {
		eventOpGraphId := this.getEventOpGraphId(event)
		isMatchFound := false
		for _, opGraphId := range filter.OpGraphIds {
			if opGraphId == eventOpGraphId {
				isMatchFound = true
				break
			}
		}
		if !isMatchFound {
			return true
		}
	}
	return false
}

// O(n) complexity (n being subscription count); thread safe
func (this *eventBus) publishToSubscribers(event model.Event) {
	this.subscriberMutex.RLock()
	defer this.subscriberMutex.RUnlock()
	for subscriber, filter := range this.subscribers {
		if !this.isEventFiltered(event, filter) {
			// use go routine; sending event may be slow or timeout
			go func(event model.Event, subscriber chan model.Event) {
				select {
				case subscriber <- event:
				case <-time.After(time.Second * 5):
					this.UnregisterSubscriber(subscriber)
				}
			}(event, subscriber)
		}
	}
}

func (this *eventBus) getEventOpGraphId(event model.Event) string {
	switch {
	case nil != event.ContainerExited:
		return event.ContainerExited.OpGraphId
	case nil != event.ContainerStarted:
		return event.ContainerStarted.OpGraphId
	case nil != event.ContainerStdErrWrittenTo:
		return event.ContainerStdErrWrittenTo.OpGraphId
	case nil != event.ContainerStdOutWrittenTo:
		return event.ContainerStdOutWrittenTo.OpGraphId
	case nil != event.OpEncounteredError:
		return event.OpEncounteredError.OpGraphId
	case nil != event.OpEnded:
		return event.OpEnded.OpGraphId
	case nil != event.OpStarted:
		return event.OpStarted.OpGraphId
	default:
		panic(fmt.Sprintf("Received unexpected event %v\n", event))
	}
}

// O(1) complexity; thread safe
func (this *eventBus) Publish(
	event model.Event,
) {
	this.pendingPublishesChannel <- event
}

// O(1) complexity; thread safe
func (this *eventBus) UnregisterSubscriber(
	eventChannel chan model.Event,
) {
	this.subscriberMutex.Lock()
	defer this.subscriberMutex.Unlock()
	if _, ok := this.subscribers[eventChannel]; ok {
		delete(this.subscribers, eventChannel)
		close(eventChannel)
	}
}
