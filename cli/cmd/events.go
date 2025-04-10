package cmd

import (
	"errors"

	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/spf13/cobra"
)

func newEventsCmd(
	cliOutput clioutput.CliOutput,
	nodeConfig *local.NodeConfig,
) *cobra.Command {
	return &cobra.Command{
		Args:  cobra.MaximumNArgs(0),
		Use:   "events",
		Short: "Stream events from an opctl node",
		Long: `If an opctl node isn't reachable, one will be started automatically. Events are delivered 
over a websocket connection. Past events are replayed when streaming starts. As new 
events occur, they are streamed in realtime.
`,
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			np, err := local.New(*nodeConfig)
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
		},
	}
}
