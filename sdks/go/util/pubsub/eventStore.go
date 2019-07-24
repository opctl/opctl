package pubsub

import (
	"context"
	"github.com/opctl/opctl/sdks/go/types"
)

type EventStore interface {
	Add(event types.Event) error
	List(
		ctx context.Context,
		filter types.EventFilter,
	) (
		<-chan types.Event,
		<-chan error,
	)
}

const sortableRFC3339Nano = "2006-01-02T15:04:05.000000000Z07:00"
