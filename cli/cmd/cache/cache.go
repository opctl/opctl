package cache

import (
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/spf13/cobra"
)

func NewCacheCmd(nodeConfig *local.NodeConfig) *cobra.Command {
	cacheCmd := cobra.Command{
		Use:   "cache",
		Short: "Manage cache dirs",
	}
	cacheCmd.AddCommand(newPruneCmd(nodeConfig))
	return &cacheCmd
}
