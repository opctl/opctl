package core

import (
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/model"
	"time"
)

func (this _core) StreamEvents() {

	// ensure node running
	nodes, err := this.nodeProvider.ListNodes()
	if nil != err {
		panic(err.Error())
	}
	if len(nodes) < 1 {
		this.nodeProvider.CreateNode()
		// sleep to let the opctl node start
		// @TODO: add exp backoff to SDK websocket client so we don't need this
		<-time.After(time.Second * 3)
	}

	eventChannel, err := this.consumeNodeApi.GetEventStream(
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
