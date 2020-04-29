package pubsub

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/opctl/opctl/sdks/go/model"
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

	// per badger README.MD#FAQ "maximizes throughput"
	runtime.GOMAXPROCS(128)

	opts := badger.DefaultOptions(eventDbDirPath).WithLogger(nil)
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
	eventChannel := make(chan model.Event, 1000)
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
			defer it.Close()
			sinceBytes := []byte(sinceTime.Format(sortableRFC3339Nano))
			for it.Seek(sinceBytes); it.Valid(); it.Next() {
				item := it.Item()
				item.Value(func(v []byte) error {
					event := model.Event{}
					if err := json.Unmarshal(v, &event); nil != err {
						return err
					}

					if !isRootOpIDExcludedByFilter(getEventRootOpID(event), filter) {
						select {
						case <-ctx.Done():
							return ctx.Err()
						case eventChannel <- event:
						}
					}
					return nil
				})
			}

			return nil
		}); nil != err {
			errChannel <- err
		}
	}()

	return eventChannel, errChannel
}
