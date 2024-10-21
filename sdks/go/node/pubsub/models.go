package pubsub

import (
	"github.com/opctl/opctl/sdks/go/model"
)

type subscriptionInfo struct {
	Filter model.EventFilter
	Done   chan struct{}
}
