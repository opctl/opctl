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
		Example: `# open web UI to current working directory
opctl ui

# open web UI to root of github.com/opspec-pkgs/_.op.create git repository tag 3.3.1
opctl ui 'github.com/opspec-pkgs/_.op.create#3.3.1'
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
