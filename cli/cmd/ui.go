package cmd

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

func newUICmd(
	cliOutput clioutput.CliOutput,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	nodeConfig *local.NodeConfig,
) *cobra.Command {
	mountRefArgName := "MOUNT_REF"

	return &cobra.Command{
		Args: cobra.MaximumNArgs(1),
		Example: `# Open the opctl web UI to the current working directory.
opctl ui

# Open the opctl web UI to the root directory of the 'github.com/opspec-pkgs/github.release.create' git repository commit tagged '3.0.0'.
opctl ui github.com/opspec-pkgs/github.release.create#3.0.0
`,
		Use: fmt.Sprintf(
			"ui [%s]",
			mountRefArgName,
		),
		Short:   "Open the opctl web UI and mount a reference",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var resolvedMount string

			np, err := local.New(*nodeConfig)
			if err != nil {
				return err
			}

			node, err := np.CreateNodeIfNotExists(ctx)
			if err != nil {
				return err
			}

			mountRefArg := ""

			if len(args) > 0 {
				mountRefArg = args[0]
			}

			if strings.HasPrefix(mountRefArg, ".") {
				// treat dot paths as regular rel paths
				resolvedMount, err = filepath.Abs(mountRefArg)
				if err != nil {
					return err
				}
			} else {
				dataResolver := dataresolver.New(
					cliParamSatisfier,
					node,
				)

				// otherwise use same resolution as run
				mountHandle, err := dataResolver.Resolve(
					ctx,
					mountRefArg,
					nil,
				)
				if err != nil {
					return err
				}

				resolvedMount = mountHandle.Ref()
			}

			err = open.Run(
				fmt.Sprintf("http://%s/?mount=%s", nodeConfig.APIListenAddress, url.QueryEscape(resolvedMount)),
			)
			if err != nil {
				return err
			}

			cliOutput.Success("Opctl web UI opened!")

			return nil
		},
	}
}
