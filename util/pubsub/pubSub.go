package pubsub

//go:generate counterfeiter -o ./fakeEventSubscriber.go --fake-name FakeEventSubscriber ./ EventSubscriber
//go:generate counterfeiter -o ./fakeEventPublisher.go --fake-name FakeEventPublisher ./ EventPublisher
//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ PubSub

import (
	"context"
	"github.com/opctl/sdk-golang/model"
	"sync"
	"time"
)

func New(
	eventStore EventStore,
) PubSub {
	return &pubSub{
		eventStore:    eventStore,
		subscriptions: map[chan model.Event]subscriptionInfo{},
	}
}

type EventPublisher interface {
	Publish(
		event model.Event,
	)
}

type EventSubscriber interface {
	// Subscribe returns a filtered event stream
	// events will be sent to the subscription until either:
	//  - ctx is canceled
	//  - returned channel is blocked for 10 seconds
	// note: method signature is based on https://medium.com/statuscode/pipeline-patterns-in-go-a37bb3a7e61d
	Subscribe(
		ctx context.Context,
		filter model.EventFilter,
	) (
		<-chan model.Event,
		<-chan error,
	)
}

type PubSub interface {
	EventPublisher
	EventSubscriber
}

type pubSub struct {
	eventStore EventStore
	// subscriptions is a map where key is a channel for the subscription & value is info about the subscription
	subscriptions      map[chan model.Event]subscriptionInfo
	subscriptionsMutex sync.RWMutex
}

func (ps *pubSub) Subscribe(
	ctx context.Context,
	filter model.EventFilter,
) (
	<-chan model.Event,
	<-chan error,
) {
	dstEventChannel := make(chan model.Event, 1000)
	dstErrChannel := make(chan error, 1)

	go func() {
		defer close(dstEventChannel)
		defer close(dstErrChannel)

		publishEventChannel := make(chan model.Event, 1000)
		defer ps.gcSubscription(publishEventChannel)

		subscriptionInfo := subscriptionInfo{
			Filter: filter,
			// Done is closed when the subscription is garbage collected
			Done: make(chan struct{}, 1),
		}

		ps.subscriptionsMutex.Lock()
		ps.subscriptions[publishEventChannel] = subscriptionInfo
		ps.subscriptionsMutex.Unlock()

		// old events
		eventStoreEventChannel, eventStoreErrChannel := ps.eventStore.List(ctx, filter)
		for event := range eventStoreEventChannel {
			select {
			case dstEventChannel <- event:
			case <-ctx.Done():
				return
			case <-time.After(time.Second * 10):
				// evict the subscriber (they didn't accept the event within 10 seconds)
				return
			}
		}

		if err := <-eventStoreErrChannel; nil != err {
			dstErrChannel <- err
			return
		}

		// new events
		for event := range publishEventChannel {
			select {
			case dstEventChannel <- event:
			case <-ctx.Done():
				return
			case <-time.After(time.Second * 10):
				// evict the subscriber (they didn't accept the event within 10 seconds)
				return
			}
		}
	}()

	return dstEventChannel, dstErrChannel
}

func (ps *pubSub) gcSubscription(
	channel chan model.Event,
) {
	ps.subscriptionsMutex.RLock()
	close(ps.subscriptions[channel].Done)
	ps.subscriptionsMutex.RUnlock()

	ps.subscriptionsMutex.Lock()
	delete(ps.subscriptions, channel)
	ps.subscriptionsMutex.Unlock()
}

// O(n) complexity (n being number of existing subscriptions); thread safe
func (ps *pubSub) Publish(
	event model.Event,
) {
	ps.eventStore.Add(event)

	ps.subscriptionsMutex.RLock()
	defer ps.subscriptionsMutex.RUnlock()

	for publishEventChannel, subscriptionInfo := range ps.subscriptions {

		RootOpID := getEventRootOpID(event)
		if !isRootOpIDExcludedByFilter(RootOpID, subscriptionInfo.Filter) {

			// use go routine because this publishEventChannel could be blocked
			// for valid reasons such as replaying events from event store.
			//
			// In such a case, we don't want to hold up delivery to any
			// other subscriptions
			go ps.publishToSubscription(publishEventChannel, subscriptionInfo, event)

		}

	}
}

/**
publishToSubscription publishes event to subscription
*/
func (ps *pubSub) publishToSubscription(
	subscriptionChannel chan model.Event,
	subscriptionInfo subscriptionInfo,
	event model.Event,
) {
	select {
	case <-subscriptionInfo.Done:
	case subscriptionChannel <- event:
	}
}
