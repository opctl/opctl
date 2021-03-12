package main

import (
	"context"
	"errors"

	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/sdks/go/model"
)

// events implements "events" command
func events(
	ctx context.Context,
	cliOutput clioutput.CliOutput,
	nodeProvider nodeprovider.NodeProvider,
) error {
	node, err := nodeProvider.CreateNodeIfNotExists(ctx)
	if err != nil {
		return err
	}

	eventChannel, err := node.GetEventStream(
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

		cliOutput.Event(&event)
	}
	return nil
}
