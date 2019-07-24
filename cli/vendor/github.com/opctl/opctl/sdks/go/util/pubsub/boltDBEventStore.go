package pubsub

import (
	"context"
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/opctl/opctl/sdks/go/types"
	"os"
	"path"
	"sync"
	"time"
)

/**
NewBoltDBEventStore returns an EventStore implementation leveraging [Bolt DB](https://github.com/boltdb/bolt)
*/
func NewBoltDBEventStore(
	eventDbFilePath string,
) EventStore {
	err := os.MkdirAll(path.Dir(eventDbFilePath), 0700)
	if nil != err {
		panic(err)
	}

	db, err := bolt.Open(eventDbFilePath, 0644, nil)
	if nil != err {
		panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("events"))
		return err
	})
	if nil != err {
		panic(err)
	}

	return &boltDBEventStore{
		db: db,
	}
}

type boltDBEventStore struct {
	db               *bolt.DB
	eventsByRootOpID map[string][]*types.Event
	eventsMutex      sync.RWMutex
}

// O(1); threadsafe
func (this *boltDBEventStore) Add(event types.Event) error {

	return this.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("events"))

		encodedEvent, err := json.Marshal(event)
		if nil != err {
			return err
		}

		return bucket.Put([]byte(event.Timestamp.Format(sortableRFC3339Nano)), encodedEvent)
	})
}

// O(n) (n being number of subscriptions that exist); threadsafe
func (this *boltDBEventStore) List(ctx context.Context,
	filter types.EventFilter,
) (<-chan types.Event, <-chan error) {
	eventChannel := make(chan types.Event, 1000)
	errChannel := make(chan error, 1)

	go func() {
		defer close(eventChannel)
		defer close(errChannel)

		if err := this.db.View(func(tx *bolt.Tx) error {
			cursor := tx.Bucket([]byte("events")).Cursor()

			sinceTime := new(time.Time)
			if nil != filter.Since {
				sinceTime = filter.Since
			}

			sinceBytes := []byte(sinceTime.Format(sortableRFC3339Nano))
			for k, v := cursor.Seek(sinceBytes); k != nil; k, v = cursor.Next() {
				event := types.Event{}
				err := json.Unmarshal(v, &event)
				if nil != err {
					return err
				}

				if !isRootOpIDExcludedByFilter(getEventRootOpID(event), filter) {
					eventChannel <- event
				}
			}

			return nil
		}); nil != err {
			errChannel <- err
		}
	}()

	return eventChannel, errChannel
}
