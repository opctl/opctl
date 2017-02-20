package pubsub

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
	"sync"
)

// interface for event storage
type eventRepo interface {
	Add(event *model.Event)
	List(filter *model.EventFilter) []*model.Event
}

func newEventRepo() eventRepo {
	return &_eventRepo{
		eventsByOgid: make(map[string][]*model.Event),
		eventsMutex:  sync.RWMutex{},
	}
}

type _eventRepo struct {
	eventsByOgid map[string][]*model.Event
	eventsMutex  sync.RWMutex
}

// O(1); threadsafe
func (this *_eventRepo) Add(event *model.Event) {
	this.eventsMutex.Lock()
	defer this.eventsMutex.Unlock()

	opGraphId := getEventOpGraphId(event)
	this.eventsByOgid[opGraphId] = append(this.eventsByOgid[opGraphId], event)
}

// O(n) (n being number of subscriptions that exist); threadsafe
func (this *_eventRepo) List(filter *model.EventFilter) []*model.Event {
	this.eventsMutex.RLock()
	defer this.eventsMutex.RUnlock()
	result := []*model.Event{}
	for opGraphId, events := range this.eventsByOgid {
		if !isOgidExcludedByFilter(opGraphId, filter) {
			result = append(result, events...)
		}
	}
	// @TODO: sort
	return result
}
