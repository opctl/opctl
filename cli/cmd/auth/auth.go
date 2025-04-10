package auth

import (
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/spf13/cobra"
)

func NewAuthCmd(
	nodeConfig *local.NodeConfig,
) *cobra.Command {
	authCmd := cobra.Command{
		Use:   "auth",
		Short: "Manage auth",
	}

	authCmd.AddCommand(
		newAddCmd(
			nodeConfig,
		),
	)

	return &authCmd
}
