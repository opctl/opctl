package op

import (
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	nodePkg "github.com/opctl/opctl/sdks/go/node"
	"github.com/spf13/cobra"
)

var node nodePkg.Node

func NewOpCmd(
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	nodeConfig *local.NodeConfig,
) *cobra.Command {
	opCmd := cobra.Command{
		Use:   "op",
		Short: "Manage ops",
	}

	opCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		np, err := local.New(*nodeConfig)
		if err != nil {
			return err
		}

		node, err = np.CreateNodeIfNotExists(
			cmd.Context(),
		)
		return err

	}

	opCmd.AddCommand(
		newCreateCmd(),
	)
	opCmd.AddCommand(
		newInstallCmd(
			cliParamSatisfier,
		),
	)
	opCmd.AddCommand(
		newKillCmd(),
	)
	opCmd.AddCommand(
		newValidateCmd(
			cliParamSatisfier,
		),
	)

	return &opCmd
}
