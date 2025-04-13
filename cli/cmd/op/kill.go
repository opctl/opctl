package op

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/spf13/cobra"
)

func newKillCmd() *cobra.Command {
	opIDArgName := "OP_ID"

	killCmd := cobra.Command{
		Args: cobra.ExactArgs(1),
		Use: fmt.Sprintf(
			"kill %s",
			opIDArgName,
		),
		Short: "Kill an op",
		RunE: func(cmd *cobra.Command, args []string) error {
			opID := args[0]

			return node.KillOp(
				cmd.Context(),
				model.KillOpReq{
					OpID:       opID,
					RootCallID: opID,
				},
			)
		},
	}

	return &killCmd
}
