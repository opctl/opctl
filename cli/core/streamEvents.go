package core

import "github.com/opspec-io/sdk-golang/pkg/model"

func (this _core) StreamEvents() {

	eventChannel, err := this.engineClient.GetEventStream(
		&model.GetEventStreamReq{},
	)
	if nil != err {
		this.exiter.Exit(ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	for {

		event, isEventChannelOpen := <-eventChannel
		if !isEventChannelOpen {
			this.exiter.Exit(ExitReq{Message: "Event channel closed unexpectedly", Code: 1})
			return // support fake exiter
		}

		this.output.Event(&event)
	}
}
