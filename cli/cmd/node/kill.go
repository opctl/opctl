package node

import (
	"github.com/opctl/opctl/cli/internal/euid0"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/dns"
	"github.com/spf13/cobra"
)

func newKillCmd(
	containerRuntime *containerruntime.ContainerRuntime,
	nodeConfig *local.NodeConfig,
) *cobra.Command {
	killCmd := cobra.Command{
		Args:  cobra.ExactArgs(0),
		Use:   "kill",
		Short: "Kill an opctl node and any running operations",
		Long:  "Killing a node is non destructive. All node data including auth, caches, and operation state will be retained.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := euid0.Ensure(); err != nil {
				return err
			}

			np, err := local.New(*nodeConfig)
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

			return (*containerRuntime).Kill(ctx)
		},
	}

	return &killCmd
}
