package pubsub

import (
	"github.com/opspec-io/sdk-golang/model"
)

type subscriptionInfo struct {
	Filter model.EventFilter
	Done   chan struct{}
}
