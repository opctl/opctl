package pubsub

import (
	"encoding/json"
	"github.com/dgraph-io/badger"
	"github.com/opspec-io/sdk-golang/model"
	"log"
	"os"
	"path"
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

	return &eventRepo{
		db:               db,
		eventsByRootOpId: make(map[string][]*model.Event),
	}
}

type eventRepo struct {
	db               *badger.DB
	eventsByRootOpId map[string][]*model.Event
}

const sortableRFC3339Nano = "2006-01-02T15:04:05.000000000Z07:00"

// O(1); threadsafe
func (this *eventRepo) Add(event *model.Event) {

	// @TODO: handle errors
	this.db.Update(func(txn *badger.Txn) error {

		encodedEvent, err := json.Marshal(event)
		if nil != err {
			return err
		}

		return txn.Set([]byte(event.Timestamp.Format(sortableRFC3339Nano)), encodedEvent)
	})
}

// O(n) (n being number of subscriptions that exist); threadsafe
func (this *eventRepo) List(filter *model.EventFilter) []*model.Event {
	var result []*model.Event

	// @TODO: handle errors
	this.db.View(func(txn *badger.Txn) error {

		sinceTime := new(time.Time)
		if nil != filter && nil != filter.Since {
			sinceTime = filter.Since
		}

		it := txn.NewIterator(badger.DefaultIteratorOptions)
		sinceBytes := []byte(sinceTime.Format(sortableRFC3339Nano))
		for it.Seek(sinceBytes); it.Valid(); it.Next() {
			value, err := it.Item().Value()
			if nil != err {
				return err
			}

			event := &model.Event{}
			if err := json.Unmarshal(value, event); nil != err {
				return err
			}

			if !isRootOpIdExcludedByFilter(getEventRootOpId(event), filter) {
				result = append(result, event)
			}
		}

		return nil
	})

	// @TODO: sort
	return result
}
