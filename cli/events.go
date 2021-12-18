package main

import (
	"context"
	"errors"

	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/model"
)

// events implements "events" command
func events(
	ctx context.Context,
	cliOutput clioutput.CliOutput,
	nodeConfig local.NodeConfig,
) error {
	np, err := local.New(nodeConfig)
	if err != nil {
		return err
	}

	node, err := np.CreateNodeIfNotExists(ctx)
	if err != nil {
		return err
	}

	eventChannel, err := node.GetEventStream(
		ctx,
		&model.GetEventStreamReq{},
	)
	if err != nil {
		return err
	}

	for {
		event, isEventChannelOpen := <-eventChannel
		if !isEventChannelOpen {
			return errors.New("connection to event stream lost")
		}

		cliOutput.Event(&event)
	}
}
