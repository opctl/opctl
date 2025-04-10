package cmd

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/euid0"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

func newSelfUpdateCmd(
	cliOutput clioutput.CliOutput,
	nodeConfig *local.NodeConfig,
) *cobra.Command {
	return &cobra.Command{
		Args:  cobra.MaximumNArgs(0),
		Use:   "self-update",
		Short: "Update opctl",
		Long: `If a node is running, it will be killed in order to apply the update.
`,
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if err := euid0.Ensure(); err != nil {
				return err
			}

			v := semver.MustParse(version)
			latest, err := selfupdate.UpdateSelf(v, "opctl/opctl")
			if err != nil {
				return err
			}

			if latest.Version.Equals(v) {
				cliOutput.Success("No update available, already at the latest version!")
				return nil
			}

			// kill local node to ensure outdated version not left running
			// @TODO start node maintaining previous user
			np, err := local.New(*nodeConfig)
			if err != nil {
				return err
			}

			err = np.KillNodeIfExists(
				ctx,
			)
			if err != nil {
				return fmt.Errorf("unable to kill running node; run `node kill` to complete the update: %w", err)
			}

			cliOutput.Success(fmt.Sprintf("Updated to new version: %s!", latest.Version))

			return nil
		},
	}
}
