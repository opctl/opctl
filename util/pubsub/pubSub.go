package pubsub

//go:generate counterfeiter -o ./fakeEventSubscriber.go --fake-name FakeEventSubscriber ./ EventSubscriber
//go:generate counterfeiter -o ./fakeEventPublisher.go --fake-name FakeEventPublisher ./ EventPublisher
//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ PubSub

import (
	"context"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/uniquestring"
	"sync"
)

func New(
	eventStore EventStore,
) PubSub {
	return &pubSub{
		eventStore:          eventStore,
		subscriptions:       map[string]subscription{},
		subscriptionsMutex:  sync.RWMutex{},
		uniqueStringFactory: uniquestring.NewUniqueStringFactory(),
	}
}

type EventPublisher interface {
	Publish(
		event model.Event,
	)
}

type EventSubscriber interface {
	// Subscribe returns a filtered event stream
	// It is up to the caller to cancel the context or the subscription will continue receiving events.
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
	uniqueStringFactory uniquestring.UniqueStringFactory
	eventStore          EventStore
	subscriptions       map[string]subscription
	subscriptionsMutex  sync.RWMutex
}

func (ps *pubSub) Subscribe(
	ctx context.Context,
	filter model.EventFilter,
) (
	<-chan model.Event,
	<-chan error,
) {
	dstEventChannel := make(chan model.Event, 100000)
	dstErrChannel := make(chan error, 1)

	subscription := subscription{
		Filter:          filter,
		NewEventChannel: make(chan model.Event, 1000000),
		// DoneChannel is closed when the subscription is garbage collected
		DoneChannel: make(chan struct{}),
	}

	go func() {
		defer close(dstEventChannel)
		defer close(dstErrChannel)

		subscriptionId, err := ps.uniqueStringFactory.Construct()
		if nil != err {
			dstErrChannel <- err
			return
		}

		ps.subscriptionsMutex.Lock()
		ps.subscriptions[subscriptionId] = subscription
		ps.subscriptionsMutex.Unlock()
		defer ps.gcSubscription(subscriptionId)

		// old events
		srcEventChannel, srcErrChannel := ps.eventStore.List(ctx, filter)
		for event := range srcEventChannel {
			select {
			case <-ctx.Done():
				return
			case dstEventChannel <- event:
			}
		}

		if err := <-srcErrChannel; nil != err {
			dstErrChannel <- err
			return
		}

		// new events
		for event := range subscription.NewEventChannel {
			select {
			case <-ctx.Done():
				return
			case dstEventChannel <- event:
			}
		}
	}()

	return dstEventChannel, dstErrChannel
}

// gcSubscription garbage collects subscription w/ subscriptionId
func (ps *pubSub) gcSubscription(
	subscriptionId string,
) {
	close(ps.subscriptions[subscriptionId].DoneChannel)
	ps.subscriptionsMutex.Lock()
	delete(ps.subscriptions, subscriptionId)
	ps.subscriptionsMutex.Unlock()
}

// O(n) complexity (n being number of existing subscriptions); thread safe
func (ps *pubSub) Publish(
	event model.Event,
) {
	ps.eventStore.Add(event)

	go func() {
		ps.subscriptionsMutex.RLock()
		defer ps.subscriptionsMutex.RUnlock()
		for _, subscription := range ps.subscriptions {
			RootOpId := getEventRootOpId(event)
			if !isRootOpIdExcludedByFilter(RootOpId, subscription.Filter) {
				select {
				case <-subscription.DoneChannel:
					close(subscription.NewEventChannel)
				case subscription.NewEventChannel <- event:
				}
			}
		}
	}()
}
