package pubsub

import (
	"github.com/opctl/sdk-golang/model"
)

type subscriptionInfo struct {
	Filter model.EventFilter
	Done   chan struct{}
}
