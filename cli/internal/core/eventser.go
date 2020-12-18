package core

import (
	"context"
	"errors"

	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/sdks/go/model"
)

// Eventser exposes the "events" command
type Eventser interface {
	Events(
		ctx context.Context,
	) error
}

// newEventser returns an initialized "events" command
func newEventser(
	cliOutput clioutput.CliOutput,
	nodeProvider nodeprovider.NodeProvider,
) Eventser {
	return _eventser{
		cliOutput:    cliOutput,
		nodeProvider: nodeProvider,
	}
}

type _eventser struct {
	cliOutput    clioutput.CliOutput
	nodeProvider nodeprovider.NodeProvider
}

func (ivkr _eventser) Events(
	ctx context.Context,
) error {
	nodeHandle, err := ivkr.nodeProvider.CreateNodeIfNotExists()
	if nil != err {
		return err
	}

	eventChannel, err := nodeHandle.APIClient().GetEventStream(
		ctx,
		&model.GetEventStreamReq{},
	)
	if nil != err {
		return err
	}

	for {
		event, isEventChannelOpen := <-eventChannel
		if !isEventChannelOpen {
			return errors.New("Connection to event stream lost")
		}

		ivkr.cliOutput.Event(&event)
	}
}
