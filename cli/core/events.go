package core

import (
	"context"
	"github.com/opctl/opctl/sdk/go/model"
	"github.com/opctl/opctl/util/cliexiter"
)

func (this _core) Events(
	ctx context.Context,
) {

	// ensure node reachable
	this.nodeReachabilityEnsurer.EnsureNodeReachable()

	eventChannel, err := this.apiClient.GetEventStream(
		ctx,
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
