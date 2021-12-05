package pubsub

import (
	"context"
	"encoding/json"
	"time"

	"github.com/opctl/opctl/sdks/go/model"

	"github.com/dgraph-io/badger/v3"
)

const sortableRFC3339Nano = "2006-01-02T15:04:05.000000000Z07:00"

//newEventStore returns an EventStore implementation leveraging [Badger DB](https://github.com/dgraph-io/badger)
func newEventStore(
	db *badger.DB,
) eventStore {
	return eventStore{
		eventsByTimestampKeyPrefix: "eventsByTimestamp_",
		db:                         db,
	}
}

type eventStore struct {
	eventsByTimestampKeyPrefix string
	db                         *badger.DB
}

// O(1); threadsafe
func (es *eventStore) Add(
	event model.Event,
) error {
	return es.db.Update(func(txn *badger.Txn) error {

		encodedEvent, err := json.Marshal(event)
		if err != nil {
			return err
		}

		return txn.Set(
			[]byte(es.eventsByTimestampKeyPrefix+event.Timestamp.Format(sortableRFC3339Nano)),
			encodedEvent,
		)
	})
}

// O(n) (n being number of events that exist); threadsafe
func (es eventStore) List(
	ctx context.Context,
	filter model.EventFilter,
	eventChannel chan model.Event,
) error {
	return es.db.View(func(txn *badger.Txn) error {
		sinceTime := new(time.Time)
		if filter.Since != nil {
			sinceTime = filter.Since
		}

		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		sinceBytes := []byte(es.eventsByTimestampKeyPrefix + sinceTime.Format(sortableRFC3339Nano))
		for it.Seek(sinceBytes); it.ValidForPrefix([]byte(es.eventsByTimestampKeyPrefix)); it.Next() {
			if err := it.Item().Value(func(v []byte) error {
				event := model.Event{}
				if err := json.Unmarshal(v, &event); err != nil {
					return err
				}

				if !isRootCallIDExcludedByFilter(getEventRootCallID(event), filter) {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case eventChannel <- event:
						return ctx.Err()
					}
				}
				return nil
			}); nil != err {
				return err
			}
		}

		return nil
	})
}
