package pubsub

import (
	"github.com/opspec-io/sdk-golang/model"
)

type subscription struct {
	Filter          model.EventFilter
	NewEventChannel chan model.Event
	DoneChannel     chan struct{}
}
