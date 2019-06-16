package pubsub

import (
	"github.com/opctl/opctl/sdk/go/model"
)

type subscriptionInfo struct {
	Filter model.EventFilter
	Done   chan struct{}
}
