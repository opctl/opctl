package pubsub

import (
	"context"
	"encoding/json"
	"time"

	"github.com/opctl/opctl/sdks/go/model"

	"github.com/dgraph-io/badger/v2"
)

// EventStore stores events outside the context of a subscription.
// It allows inspecting events that have happened before a subscription is created.
type EventStore interface {
	Add(event model.Event) error
	List(
		ctx context.Context,
		filter model.EventFilter,
		eventChannel chan model.Event,
	) error
}

const sortableRFC3339Nano = "2006-01-02T15:04:05.000000000Z07:00"

// newEventStore returns an EventStore implementation leveraging [Badger DB](https://github.com/dgraph-io/badger)
func newEventStore(
	db *badger.DB,
) EventStore {
	return &_eventStore{
		eventsByTimestampKeyPrefix: "eventsByTimestamp_",
		db:                         db,
	}
}

type _eventStore struct {
	eventsByTimestampKeyPrefix string
	db                         *badger.DB
}

// O(1); threadsafe
func (es *_eventStore) Add(
	event model.Event,
) error {
	return es.db.Update(func(txn *badger.Txn) error {

		encodedEvent, err := json.Marshal(event)
		if nil != err {
			return err
		}

		return txn.Set(
			[]byte(es.eventsByTimestampKeyPrefix+event.Timestamp.Format(sortableRFC3339Nano)),
			encodedEvent,
		)
	})
}

// List sends events that occurred to an event channel. It's intended to be used
// to replay events that happened before a subscription channel was created.
//
// O(n) (n being number of events that exist); threadsafe
func (es _eventStore) List(
	ctx context.Context,
	filter model.EventFilter,
	eventChannel chan model.Event,
) error {
	if err := es.db.View(func(txn *badger.Txn) error {
		sinceTime := new(time.Time)
		if nil != filter.Since {
			sinceTime = filter.Since
		}

		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		sinceBytes := []byte(es.eventsByTimestampKeyPrefix + sinceTime.Format(sortableRFC3339Nano))
		for it.Seek(sinceBytes); it.ValidForPrefix([]byte(es.eventsByTimestampKeyPrefix)); it.Next() {
			if err := it.Item().Value(func(v []byte) error {
				var event model.Event
				if err := json.Unmarshal(v, &event); nil != err {
					return err
				}

				if !isRootCallIDExcludedByFilter(getEventRootCallID(event), filter) {
					select {
					case <-ctx.Done():
					case eventChannel <- event:
						return ctx.Err()
					}
				}
				return nil
			}); err != nil {
				return err
			}
		}

		return nil
	}); nil != err {
		return err
	}

	return nil
}
