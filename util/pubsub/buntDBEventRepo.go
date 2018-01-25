package pubsub

import (
	"context"
	"encoding/json"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/tidwall/buntdb"
	"os"
	"path"
	"sync"
	"time"
)

func NewBuntDBEventRepo(
	eventDbFilePath string,
) EventRepo {
	err := os.MkdirAll(path.Dir(eventDbFilePath), 0700)
	if nil != err {
		panic(err)
	}

	db, err := buntdb.Open(eventDbFilePath)
	if nil != err {
		panic(err)
	}

	return &buntDBEventRepo{
		db: db,
	}
}

type buntDBEventRepo struct {
	db               *buntdb.DB
	eventsByRootOpId map[string][]*model.Event
	eventsMutex      sync.RWMutex
}

// O(1); threadsafe
func (this *buntDBEventRepo) Add(event model.Event) error {

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
func (this *buntDBEventRepo) List(ctx context.Context,
	filter model.EventFilter,
) (<-chan model.Event, <-chan error) {
	eventChannel := make(chan model.Event, 100000)
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

				if !isRootOpIdExcludedByFilter(getEventRootOpId(event), filter) {
					eventChannel <- event
				}
				return true
			})
		}); nil != err {
			errChannel <- err
		}
	}()

	return eventChannel, errChannel
}
