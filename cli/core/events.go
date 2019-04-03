package core

import (
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/sdk-golang/model"
)

func (this _core) Events() {

	// ensure node reachable
	this.nodeReachabilityEnsurer.EnsureNodeReachable()

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
			this.cliExiter.Exit(
				cliexiter.ExitReq{
					Message: "Connection to event stream lost",
					Code:    1,
				},
			)
			return // support fake exiter
		}

		this.cliOutput.Event(&event)
	}
}
