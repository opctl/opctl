package pubsub

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/dgraph-io/badger/v4"
	"github.com/opctl/opctl/sdks/go/model"
)

type EventPublisher interface {
	Publish(
		event model.Event,
	)
}

type EventSubscriber interface {
	// Subscribe returns a filtered event stream
	// events will be sent to the subscription until ctx is canceled.
	// note: method signature is based on https://medium.com/statuscode/pipeline-patterns-in-go-a37bb3a7e61d
	Subscribe(
		ctx context.Context,
		filter model.EventFilter,
	) (
		<-chan model.Event,
		error,
	)
}

//counterfeiter:generate -o fakes/pubSub.go . PubSub
type PubSub interface {
	EventPublisher
	EventSubscriber
}

func New(
	db *badger.DB,
) PubSub {
	return &pubSub{
		eventStore:    newEventStore(db),
		subscriptions: map[chan model.Event]subscriptionInfo{},
	}
}

type pubSub struct {
	eventStore eventStore
	// subscriptions is a map where key is a channel for the subscription & value is info about the subscription
	subscriptions      map[chan model.Event]subscriptionInfo
	subscriptionsMutex sync.RWMutex
}

func (ps *pubSub) Subscribe(
	ctx context.Context,
	filter model.EventFilter,
) (
	<-chan model.Event,
	error,
) {
	dstEventChannel := make(chan model.Event)

	go func() {
		defer func() {
			// don't let panics from any operation kill the server.
			if panic := recover(); panic != nil {
				fmt.Printf("recovered from panic: %s\n%s", panic, string(debug.Stack()))
			}
		}()

		defer close(dstEventChannel)

		publishEventChannel := make(chan model.Event)
		defer ps.closeSubscription(publishEventChannel)

		subscriptionInfo := subscriptionInfo{
			Filter: filter,
			// Done is closed when the subscription is garbage collected
			Done: make(chan struct{}, 1),
		}

		ps.subscriptionsMutex.Lock()
		ps.subscriptions[publishEventChannel] = subscriptionInfo
		ps.subscriptionsMutex.Unlock()

		// old events
		err := ps.eventStore.List(ctx, filter, dstEventChannel)
		if err != nil {
			return
		}

		// new events
		for event := range publishEventChannel {
			select {
			case <-ctx.Done():
				return
			case dstEventChannel <- event:
			}
		}
	}()

	return dstEventChannel, nil
}

func (ps *pubSub) closeSubscription(
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

		rootCallID := getEventRootCallID(event)
		if !isRootCallIDExcludedByFilter(rootCallID, subscriptionInfo.Filter) {

			// use go routine because this publishEventChannel could be blocked
			// for valid reasons such as replaying events from event store.
			//
			// In such a case, we don't want to hold up delivery to any
			// other subscriptions
			go ps.publishToSubscription(publishEventChannel, subscriptionInfo, event)

		}

	}
}

/*
*
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
