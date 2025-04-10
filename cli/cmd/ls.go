package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/spf13/cobra"
)

func newLsCmd(
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	nodeConfig *local.NodeConfig,
) *cobra.Command {
	dirRefArgName := "DIR_REF"

	return &cobra.Command{
		Args: cobra.MaximumNArgs(1),
		Example: `# list ops at ./.opspec
opctl ls  

# list ops at /absolute/path
opctl ls '/absolute/path'

# list ops at root of github.com/opctl/opctl git repository tag 0.1.24
opctl ls 'github.com/opctl/opctl#0.1.24'
`,
		Use: fmt.Sprintf(
			"ls [%s]",
			dirRefArgName,
		),
		Short:   "List operations",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			np, err := local.New(*nodeConfig)
			if err != nil {
				return err
			}

			node, err := np.CreateNodeIfNotExists(ctx)
			if err != nil {
				return err
			}

			dataResolver := dataresolver.New(
				cliParamSatisfier,
				node,
			)

			_tabWriter := new(tabwriter.Writer)
			defer _tabWriter.Flush()
			_tabWriter.Init(os.Stdout, 0, 8, 1, '\t', 0)

			fmt.Fprintln(_tabWriter, "REF\tDESCRIPTION")

			dirRef := ""

			if len(args) > 0 {
				dirRef = args[0]
			}

			dirHandle, err := dataResolver.Resolve(
				ctx,
				dirRef,
				nil,
			)
			if err != nil {
				return err
			}

			opsByPath, err := opspec.List(
				ctx,
				dirHandle,
			)
			if err != nil {
				return err
			}

			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			for path, op := range opsByPath {
				opRef := filepath.Join(dirHandle.Ref(), path)
				if filepath.IsAbs(opRef) {
					// make absolute paths relative
					relOpRef, err := filepath.Rel(cwd, opRef)
					if err != nil {
						return err
					}

					opRef = strings.TrimPrefix(relOpRef, ".opspec/")
				}

				description := op.Description
				if strings.TrimSpace(description) == "" {
					description = "-"
				}

				scanner := bufio.NewScanner(strings.NewReader(description))
				if scanner.Scan() {
					// first line of description, add the op ref
					fmt.Fprintf(_tabWriter, "%v\t%v", opRef, scanner.Text())
				}
				for scanner.Scan() {
					// subsequent lines, don't add the op ref but let the description span multiple lines
					fmt.Fprintf(_tabWriter, "\n\t%v", scanner.Text())
				}
				fmt.Fprintln(_tabWriter)
			}

			return nil
		},
	}
}
