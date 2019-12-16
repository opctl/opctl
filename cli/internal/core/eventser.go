package core

import (
	"context"

	"github.com/opctl/opctl/cli/internal/apireachabilityensurer"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

// Eventser exposes the "events" command
type Eventser interface {
	Events(
		ctx context.Context,
	)
}

// newEventser returns an initialized "events" command
func newEventser(
	apiClient client.Client,
	cliExiter cliexiter.CliExiter,
	cliOutput clioutput.CliOutput,
) Eventser {
	return _eventser{
		apiClient:              apiClient,
		apiReachabilityEnsurer: apireachabilityensurer.New(cliExiter),
		cliExiter:              cliExiter,
		cliOutput:              cliOutput,
	}
}

type _eventser struct {
	apiClient              client.Client
	apiReachabilityEnsurer apireachabilityensurer.APIReachabilityEnsurer
	cliExiter              cliexiter.CliExiter
	cliOutput              clioutput.CliOutput
}

func (ivkr _eventser) Events(
	ctx context.Context,
) {

	// ensure node reachable
	ivkr.apiReachabilityEnsurer.Ensure()

	eventChannel, err := ivkr.apiClient.GetEventStream(
		ctx,
		&model.GetEventStreamReq{},
	)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	for {

		event, isEventChannelOpen := <-eventChannel
		if !isEventChannelOpen {
			ivkr.cliExiter.Exit(
				cliexiter.ExitReq{
					Message: "Connection to event stream lost",
					Code:    1,
				},
			)
			return // support fake exiter
		}

		ivkr.cliOutput.Event(&event)
	}
}
