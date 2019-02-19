package pubsub

import (
	"context"
	"encoding/json"
	"github.com/opctl/sdk-golang/model"
	"github.com/tidwall/buntdb"
	"os"
	"path"
	"sync"
	"time"
)

/**
NewBuntDBEventStore returns an EventStore implementation leveraging [Bunt DB](https://github.com/tidwall/buntdb)
*/
func NewBuntDBEventStore(
	eventDbFilePath string,
) EventStore {
	err := os.MkdirAll(path.Dir(eventDbFilePath), 0700)
	if nil != err {
		panic(err)
	}

	db, err := buntdb.Open(eventDbFilePath)
	if nil != err {
		panic(err)
	}

	return &buntDBEventStore{
		db: db,
	}
}

type buntDBEventStore struct {
	db               *buntdb.DB
	eventsByRootOpID map[string][]*model.Event
	eventsMutex      sync.RWMutex
}

// O(1); threadsafe
func (this *buntDBEventStore) Add(event model.Event) error {

	return this.db.Update(func(tx *buntdb.Tx) error {

		encodedEvent, err := json.Marshal(event)
		if nil != err {
			return err
		}

		_, _, err = tx.Set(event.Timestamp.Format(sortableRFC3339Nano), string(encodedEvent), nil)
		return err
	})
}

// O(n) (n being number of subscriptions that exist); threadsafe
func (this *buntDBEventStore) List(ctx context.Context,
	filter model.EventFilter,
) (<-chan model.Event, <-chan error) {
	eventChannel := make(chan model.Event, 1000)
	errChannel := make(chan error, 1)

	go func() {
		defer close(eventChannel)
		defer close(errChannel)

		if err := this.db.View(func(tx *buntdb.Tx) error {

			sinceTime := new(time.Time)
			if nil != filter.Since {
				sinceTime = filter.Since
			}

			return tx.AscendGreaterOrEqual("", sinceTime.Format(sortableRFC3339Nano), func(key, value string) bool {
				event := model.Event{}
				err := json.Unmarshal([]byte(value), &event)
				if nil != err {
					return false
				}

				if isRootOpIDExcludedByFilter(getEventRootOpID(event), filter) {
					return true
				}

				select {
				case <-ctx.Done():
					return false
				case eventChannel <- event:
					return true
				}
			})
		}); nil != err {
			errChannel <- err
		}
	}()

	return eventChannel, errChannel
}
