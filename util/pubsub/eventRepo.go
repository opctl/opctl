package pubsub

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/opspec-io/sdk-golang/model"
	"os"
	"path"
	"sync"
	"time"
)

// interface for event storage
type EventRepo interface {
	Add(event *model.Event)
	List(filter *model.EventFilter) []*model.Event
}

func NewEventRepo(
	eventDbFilePath string,
) EventRepo {
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

	return &eventRepo{
		db:           db,
		eventsByRootOpId: make(map[string][]*model.Event),
		eventsMutex:  sync.RWMutex{},
	}
}

type eventRepo struct {
	db           *bolt.DB
	eventsByRootOpId map[string][]*model.Event
	eventsMutex  sync.RWMutex
}

const sortableRFC3339Nano = "2006-01-02T15:04:05.000000000Z07:00"

// O(1); threadsafe
func (this *eventRepo) Add(event *model.Event) {

	// @TODO: handle errors
	this.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("events"))

		encodedEvent, err := json.Marshal(event)
		if nil != err {
			return err
		}

		return bucket.Put([]byte(event.Timestamp.Format(sortableRFC3339Nano)), encodedEvent)
	})
}

// O(n) (n being number of subscriptions that exist); threadsafe
func (this *eventRepo) List(filter *model.EventFilter) []*model.Event {
	result := []*model.Event{}

	// @TODO: handle errors
	this.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket([]byte("events")).Cursor()

		sinceTime := new(time.Time)
		if nil != filter && nil != filter.Since {
			sinceTime = filter.Since
		}

		sinceBytes := []byte(sinceTime.Format(sortableRFC3339Nano))
		for k, v := cursor.Seek(sinceBytes); k != nil; k, v = cursor.Next() {
			event := &model.Event{}
			err := json.Unmarshal(v, event)
			if nil != err {
				return err
			}

			if !isRootOpIdExcludedByFilter(getEventRootOpId(event), filter) {
				result = append(result, event)
			}
		}

		return nil
	})

	this.eventsMutex.RLock()
	defer this.eventsMutex.RUnlock()

	// @TODO: sort
	return result
}
