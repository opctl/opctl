package op

import (
	"fmt"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	nameArgName := "NAME"
	pathFlagName := "path"
	descriptionFlagName := "description"

	pathFlag := ""
	descriptionFlag := ""

	createCmd := cobra.Command{
		Args: cobra.ExactArgs(1),
		Use: fmt.Sprintf(
			"create %s",
			nameArgName,
		),
		Short: "Create an op",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			return opspec.Create(
				filepath.Join(pathFlag, name),
				name,
				descriptionFlag,
			)
		},
	}

	createCmd.Flags().StringVarP(&pathFlag, pathFlagName, "p", model.DotOpspecDirName, "Path the op will be created at")
	createCmd.Flags().StringVarP(&descriptionFlag, descriptionFlagName, "d", "", "Description of the op")

	return &createCmd
}
