package pubsub

import (
	"context"
	"encoding/json"
	"github.com/dgraph-io/badger"
	"github.com/opspec-io/sdk-golang/model"
	"log"
	"os"
	"path"
	"time"
)

/**
NewBadgerDBEventStore returns an EventStore implementation leveraging [Badger DB](https://github.com/dgraph-io/badger)
 */
func NewBadgerDBEventStore(
	eventDbFilePath string,
) EventStore {
	eventDbDirPath := path.Dir(eventDbFilePath)
	err := os.MkdirAll(eventDbDirPath, 0700)
	if nil != err {
		panic(err)
	}

	opts := badger.DefaultOptions
	opts.Dir = eventDbDirPath
	opts.ValueDir = eventDbDirPath
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

	return &badgerDBEventStore{
		db,
	}
}

type badgerDBEventStore struct {
	db *badger.DB
}

// O(1); threadsafe
func (er *badgerDBEventStore) Add(
	event model.Event,
) error {

	return er.db.Update(func(txn *badger.Txn) error {

		encodedEvent, err := json.Marshal(event)
		if nil != err {
			return err
		}

		return txn.Set([]byte(event.Timestamp.Format(sortableRFC3339Nano)), encodedEvent)
	})
}

// O(n) (n being number of events that exist); threadsafe
func (er badgerDBEventStore) List(
	ctx context.Context,
	filter model.EventFilter,
) (<-chan model.Event, <-chan error) {
	eventChannel := make(chan model.Event, 100000)
	errChannel := make(chan error, 1)

	go func() {
		defer close(eventChannel)
		defer close(errChannel)

		if err := er.db.View(func(txn *badger.Txn) error {
			sinceTime := new(time.Time)
			if nil != filter.Since {
				sinceTime = filter.Since
			}

			it := txn.NewIterator(badger.DefaultIteratorOptions)
			sinceBytes := []byte(sinceTime.Format(sortableRFC3339Nano))
			for it.Seek(sinceBytes); it.Valid(); it.Next() {
				value, err := it.Item().Value()
				if nil != err {
					return err
				}

				event := model.Event{}
				if err := json.Unmarshal(value, &event); nil != err {
					return err
				}

				if !isRootOpIdExcludedByFilter(getEventRootOpId(event), filter) {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case eventChannel <- event:
					}
				}
			}

			return nil
		}); nil != err {
			errChannel <- err
		}
	}()

	return eventChannel, errChannel
}
