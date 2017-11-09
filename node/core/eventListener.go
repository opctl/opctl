package core

//go:generate counterfeiter -o ./fakeEventListener.go --fake-name fakeEventListener ./ eventListener

import (
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/pubsub"
)

type eventListener interface {
	Listen()
}

func newEventListener(
	opKiller opKiller,
	eventSubscriber pubsub.EventSubscriber,
) eventListener {
	return _eventListener{
		opKiller:        opKiller,
		eventSubscriber: eventSubscriber,
	}
}

type _eventListener struct {
	eventSubscriber pubsub.EventSubscriber
	opKiller        opKiller
}

func (el _eventListener) Listen() {
	go func() {
		eventChannel := make(chan *model.Event, 1)
		el.eventSubscriber.Subscribe(
			nil,
			eventChannel,
		)
		for event := range eventChannel {
			if nil != event.OpKilled {
				el.opKiller.Kill(event.OpKilled.RootOpId)
			}
		}
	}()
}
