package core

import (
	"context"

	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/sdks/go/model"
)

// Eventser exposes the "events" command
type Eventser interface {
	Events(
		ctx context.Context,
	)
}

// newEventser returns an initialized "events" command
func newEventser(
	cliExiter cliexiter.CliExiter,
	cliOutput clioutput.CliOutput,
	nodeProvider nodeprovider.NodeProvider,
) Eventser {
	return _eventser{
		cliExiter:    cliExiter,
		cliOutput:    cliOutput,
		nodeProvider: nodeProvider,
	}
}

type _eventser struct {
	cliExiter    cliexiter.CliExiter
	cliOutput    clioutput.CliOutput
	nodeProvider nodeprovider.NodeProvider
}

func (ivkr _eventser) Events(
	ctx context.Context,
) {
	nodeHandle, createNodeIfNotExistsErr := ivkr.nodeProvider.CreateNodeIfNotExists()
	if nil != createNodeIfNotExistsErr {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: createNodeIfNotExistsErr.Error(), Code: 1})
		return // support fake exiter
	}

	eventChannel, err := nodeHandle.APIClient().GetEventStream(
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
