package op

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/spf13/cobra"
)

func newInstallCmd(
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
) *cobra.Command {
	opRefArgName := "OP_REF"
	pathFlagName := "path"
	usernameFlagName := "username"
	passwordFlagName := "password"

	pathFlag := ""
	usernameFlag := ""
	passwordFlag := ""

	installCmd := cobra.Command{
		Args: cobra.ExactArgs(1),
		Example: `# Install the op defined at the root of the 'github.com/opspec-pkgs/uuid.v4.generate' 
# git repository commit tagged '1.1.0' in the '.opspec/github.com/opspec-pkgs/uuid.v4.generate#1.1.0' directory
# of the current working directory.
opctl op install github.com/opspec-pkgs/uuid.v4.generate#1.1.0
`,
		Use: fmt.Sprintf(
			"install %s",
			opRefArgName,
		),
		Short: "Install an op",
		Long: `OP_REF can be either a 'host/repo-path#tag' or 'host/repo-path#tag/path'.

If auth w/ the op source fails the CLI will (re)prompt for username &
password. In non-interactive terminals, the CLI will note that it can't prompt due to being in a
non-interactive terminal and exit with a non zero exit code.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			opRef := args[0]

			// install the whole pkg in case relative (intra pkg) refs exist
			opRefParts := strings.Split(opRef, "#")
			if len(opRefParts) == 1 {
				return fmt.Errorf("%s must be a remote reference formatted as host/path#semver", opRefArgName)
			}

			var version string
			if verAndPathParts := strings.SplitN(opRefParts[1], "/", 2); len(verAndPathParts) != 1 {
				version = verAndPathParts[0]
			} else {
				version = opRefParts[1]
			}

			dataRef := fmt.Sprintf("%s#%s", opRefParts[0], version)

			dataResolver := dataresolver.New(
				cliParamSatisfier,
				node,
			)

			var creds *model.Creds
			if usernameFlag != "" && passwordFlag != "" {
				creds = &model.Creds{
					Username: usernameFlag,
					Password: passwordFlag,
				}
			}

			opDirHandle, err := dataResolver.Resolve(
				ctx,
				dataRef,
				creds,
			)
			if err != nil {
				return err
			}

			return opspec.Install(
				ctx,
				filepath.Join(pathFlag, dataRef),
				opDirHandle,
			)
		},
	}

	installCmd.Flags().StringVarP(&pathFlag, pathFlagName, "", model.DotOpspecDirName, "Path the op will be installed at")
	installCmd.Flags().StringVarP(&usernameFlag, usernameFlagName, "u", "", "Username used to auth w/ the pkg source")
	installCmd.Flags().StringVarP(&passwordFlag, passwordFlagName, "p", "", "Password used to auth w/ the pkg source")

	return &installCmd
}
