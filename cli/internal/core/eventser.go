package core

import (
	"context"
	"errors"

	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
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
	opNode node.OpNode,
) Eventser {
	return _eventser{
		cliOutput: cliOutput,
		opNode:    opNode,
	}
}

type _eventser struct {
	cliOutput clioutput.CliOutput
	opNode    node.OpNode
}

func (ivkr _eventser) Events(
	ctx context.Context,
) error {
	eventChannel, err := ivkr.opNode.GetEventStream(
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
