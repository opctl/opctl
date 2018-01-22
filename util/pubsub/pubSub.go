package pubsub

//go:generate counterfeiter -o ./fakeEventSubscriber.go --fake-name FakeEventSubscriber ./ EventSubscriber
//go:generate counterfeiter -o ./fakeEventPublisher.go --fake-name FakeEventPublisher ./ EventPublisher
//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ PubSub

import (
	"context"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/uniquestring"
	"sync"
	"time"
)

func New(
	eventRepo EventRepo,
) PubSub {
	return &pubSub{
		eventRepo:           eventRepo,
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
	eventRepo           EventRepo
	uniqueStringFactory uniquestring.UniqueStringFactory
	subscriptions       map[string]subscription
	subscriptionsMutex  sync.RWMutex
}

// thread safe
func (ps *pubSub) Subscribe(
	ctx context.Context,
	filter model.EventFilter,
) (
	<-chan model.Event,
	<-chan error,
) {

	subscriptionId := ps.uniqueStringFactory.Construct()
	subscription := subscription{
		Filter:          filter,
		NewEventChannel: make(chan model.Event, 1000000),
		// DoneChannel is closed when the subscription is garbage collected
		DoneChannel: make(chan struct{}),
	}

	ps.subscriptionsMutex.Lock()
	ps.subscriptions[subscriptionId] = subscription
	ps.subscriptionsMutex.Unlock()

	dstEventChannel := make(chan model.Event, 100000)
	dstErrChannel := make(chan error, 1)

	go func() {
		defer close(dstEventChannel)
		defer close(dstErrChannel)

		// old events
		srcEventChannel, srcErrChannel := ps.eventRepo.List(ctx, filter)
		for event := range srcEventChannel {
			select {
			case <-ctx.Done():
				ps.gcSubscription(subscriptionId)
				return
			case err, ok := <-srcErrChannel:
				if ok {
					dstErrChannel <- err
					ps.gcSubscription(subscriptionId)
					return
				}
			case dstEventChannel <- event:
			case <-time.After(time.Second * 20):
				ps.gcSubscription(subscriptionId)
				return
			}
		}

		// new events
		for event := range subscription.NewEventChannel {
			select {
			case <-ctx.Done():
				ps.gcSubscription(subscriptionId)
				return
			case dstEventChannel <- event:
			case <-time.After(time.Second * 20):
				ps.gcSubscription(subscriptionId)
				return
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
	ps.eventRepo.Add(event)

	go func() {
		ps.subscriptionsMutex.RLock()
		for _, subscription := range ps.subscriptions {
			RootOpId := getEventRootOpId(event)
			if !isRootOpIdExcludedByFilter(RootOpId, subscription.Filter) {
				select {
				case <-subscription.DoneChannel:
				case subscription.NewEventChannel <- event:
				}
			}
		}
		ps.subscriptionsMutex.RUnlock()
	}()
}
