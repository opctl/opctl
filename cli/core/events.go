package core

import (
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

func (this _core) StreamEvents() {

	eventChannel, err := this.apiClient.GetEventStream(
		&model.GetEventStreamReq{},
	)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	for {

		event, isEventChannelOpen := <-eventChannel
		if !isEventChannelOpen {
			this.cliExiter.Exit(cliexiter.ExitReq{Message: "Event channel closed unexpectedly", Code: 1})
			return // support fake exiter
		}

		this.cliOutput.Event(&event)
	}
}
