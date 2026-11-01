package op

import (
	"fmt"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/spf13/cobra"
)

// The release build links the Go op SDK release that matches this CLI.
var sdkVersion, sdkSum, sdkGoModSum string

func newCreateCmd() *cobra.Command {
	nameArgName := "NAME"
	pathFlagName := "path"
	descriptionFlagName := "description"
	goFlagName := "go"

	pathFlag := ""
	descriptionFlag := ""
	goFlag := false

	createCmd := cobra.Command{
		Args: cobra.ExactArgs(1),
		Use: fmt.Sprintf(
			"create %s",
			nameArgName,
		),
		Short: "Create an op",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			if goFlag {
				sdk := opspec.GoSDK{Version: sdkVersion, Sum: sdkSum, GoModSum: sdkGoModSum}
				if err := opspec.CreateGo(filepath.Join(pathFlag, name), name, descriptionFlag, sdk); err != nil {
					return err
				}
				if !sdk.Published() {
					fmt.Fprintln(cmd.ErrOrStderr(), "This opctl build has no published Go op SDK. Add a replace directive for github.com/opctl/opctl/sdks/go/opspec/sdk to go.mod before running the op.")
				}
				return nil
			}

			return opspec.Create(
				filepath.Join(pathFlag, name),
				name,
				descriptionFlag,
			)
		},
	}

	createCmd.Flags().StringVarP(&pathFlag, pathFlagName, "p", model.DotOpspecDirName, "Path the op will be created at")
	createCmd.Flags().StringVarP(&descriptionFlag, descriptionFlagName, "d", "", "Description of the op")
	createCmd.Flags().BoolVar(&goFlag, goFlagName, false, "Create a Go op in its own Go module")

	return &createCmd
}
