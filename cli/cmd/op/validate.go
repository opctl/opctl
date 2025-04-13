package op

import (
	"fmt"

	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"github.com/spf13/cobra"
)

func newValidateCmd(
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
) *cobra.Command {
	opRefArgName := "OP_REF"

	validateCmd := cobra.Command{
		Args: cobra.ExactArgs(1),
		Use: fmt.Sprintf(
			"validate %s",
			opRefArgName,
		),
		Short: "Validate an op",
		Long: `OP_REF can be either a 'relative/path', '/absolute/path', 'host/repo-path#tag', or 'host/repo-path#tag/path'.

Ensures:
- op.yml exists
- op.yml is syntactically valid

If auth w/ the op source fails the CLI will (re)prompt for username &
password. In non-interactive terminals, the CLI will note that it can't prompt due to being in a
non-interactive terminal and exit with a non zero exit code.
`,
		Example: `# Validate the op defined in the '.opspec/myOp' directory of the current working directory.
opctl op validate myOp

# Validate the op defined in the root directory of the 'github.com/opspec-pkgs/slack.chat.post-message' git 
# repository commit tagged '1.1.0'
opctl op validate github.com/opspec-pkgs/slack.chat.post-message#1.1.0
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			dataResolver := dataresolver.New(
				cliParamSatisfier,
				node,
			)

			opDirHandle, err := dataResolver.Resolve(
				ctx,
				args[0],
				nil,
			)
			if err != nil {
				return err
			}

			_, err = opfile.Get(
				ctx,
				opDirHandle,
			)

			return err
		},
	}

	return &validateCmd
}
