package node

import (
	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/spf13/cobra"
)

var (
	containerRuntime = new(containerruntime.ContainerRuntime)
)

func NewNodeCmd(
	cliColorer clicolorer.CliColorer,
	nodeConfig *local.NodeConfig,
) *cobra.Command {
	nodeCmd := cobra.Command{
		Use:   "node",
		Short: "Manage nodes",
	}

	nodeCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		var err error
		*containerRuntime, err = getContainerRuntime(
			cmd.Context(),
			*nodeConfig,
		)

		return err
	}

	nodeCmd.AddCommand(
		newCreateCmd(
			cliColorer,
			containerRuntime,
			nodeConfig,
		),
	)
	nodeCmd.AddCommand(
		newDeleteCmd(
			containerRuntime,
			nodeConfig,
		),
	)
	nodeCmd.AddCommand(
		newKillCmd(
			containerRuntime,
			nodeConfig,
		),
	)

	return &nodeCmd
}
