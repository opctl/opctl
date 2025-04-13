package node

import (
	"os"

	"github.com/opctl/opctl/cli/internal/euid0"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/dns"
	"github.com/spf13/cobra"
)

func newDeleteCmd(
	containerRuntime *containerruntime.ContainerRuntime,
	nodeConfig *local.NodeConfig,
) *cobra.Command {
	deleteCmd := cobra.Command{
		Args:  cobra.ExactArgs(0),
		Use:   "delete",
		Short: "Delete an opctl node",
		Long:  "Deleting a node is destructive! All node data including auth, caches, and operation state will be permanently removed.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := euid0.Ensure(); err != nil {
				return err
			}

			np, err := local.New(*nodeConfig)
			if err != nil {
				return err
			}

			err = (*containerRuntime).Delete(ctx)
			if err != nil {
				return err
			}

			if err := np.KillNodeIfExists(
				ctx,
			); err != nil {
				return err
			}

			if err := dns.DeleteResolverCfgs(
				ctx,
			); err != nil {
				return err
			}

			return os.RemoveAll(nodeConfig.DataDir)
		},
	}

	return &deleteCmd
}
