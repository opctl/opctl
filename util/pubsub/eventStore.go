package pubsub

import (
	"context"
	"github.com/opctl/sdk-golang/model"
)

type EventStore interface {
	Add(event model.Event) error
	List(
		ctx context.Context,
		filter model.EventFilter,
	) (
		<-chan model.Event,
		<-chan error,
	)
}

const sortableRFC3339Nano = "2006-01-02T15:04:05.000000000Z07:00"
