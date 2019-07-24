package pubsub

import (
	"github.com/opctl/opctl/sdks/go/types"
)

type subscriptionInfo struct {
	Filter types.EventFilter
	Done   chan struct{}
}
