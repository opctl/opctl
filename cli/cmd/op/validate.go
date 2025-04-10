package op

import (
	"fmt"

	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/opspec"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			opRef := args[0]

			dataResolver := dataresolver.New(
				cliParamSatisfier,
				node,
			)

			opDirHandle, err := dataResolver.Resolve(
				ctx,
				opRef,
				nil,
			)
			if err != nil {
				return err
			}

			return opspec.Validate(
				ctx,
				*opDirHandle.Path(),
			)
		},
	}

	return &validateCmd
}
