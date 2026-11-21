package pubsub

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
)

type orderedSubscription struct {
	filter model.EventFilter
	events []model.Event // guarded by orderedMu
	ready  chan struct{}
}

func (ps *pubSub) subscribeOrdered(ctx context.Context, filter model.EventFilter) <-chan model.Event {
	stream := make(chan model.Event)
	sub := &orderedSubscription{
		filter: filter,
		ready:  make(chan struct{}, 1),
	}

	// Register the live subscription and take a consistent replay snapshot
	// together so a terminal chunk cannot appear in both streams.
	ps.orderedMu.Lock()
	snapshot := ps.eventStore.db.NewTransaction(false)
	ps.ordered[sub] = struct{}{}
	ps.orderedMu.Unlock()

	go func() {
		defer close(stream)
		defer func() {
			ps.orderedMu.Lock()
			delete(ps.ordered, sub)
			ps.orderedMu.Unlock()
		}()

		err := ps.eventStore.list(snapshot, ctx, filter, stream)
		snapshot.Discard()
		if err != nil {
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-sub.ready:
			}
			ps.orderedMu.Lock()
			batch := sub.events
			sub.events = nil
			ps.orderedMu.Unlock()

			for i, event := range batch {
				select {
				case <-ctx.Done():
					return
				case stream <- event:
					batch[i] = model.Event{}
				}
			}
		}
	}()

	return stream
}

func (ps *pubSub) storeAndPublishOrdered(event model.Event) {
	rootID := getEventRootCallID(event)
	ps.orderedMu.RLock()
	if !ps.hasOrderedSubscriber(rootID) {
		// Protect the replay snapshot without serializing writes that have no
		// matching live subscriber.
		ps.eventStore.Add(event)
		ps.orderedMu.RUnlock()
		return
	}
	ps.orderedMu.RUnlock()

	ps.orderedMu.Lock()
	defer ps.orderedMu.Unlock()
	// Store and enqueue atomically with respect to subscription registration.
	// A replay snapshot taken between them could deliver the same event twice.
	ps.eventStore.Add(event)

	for sub := range ps.ordered {
		if isRootCallIDExcludedByFilter(rootID, sub.filter) {
			continue
		}
		sub.events = append(sub.events, event)
		select {
		case sub.ready <- struct{}{}:
		default:
		}
	}
}

// hasOrderedSubscriber reports whether any live ordered subscription matches rootID.
// The caller holds orderedMu.
func (ps *pubSub) hasOrderedSubscriber(rootID string) bool {
	for sub := range ps.ordered {
		if !isRootCallIDExcludedByFilter(rootID, sub.filter) {
			return true
		}
	}
	return false
}
