package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/appdataspec/sdk-golang/appdatapath"
	"github.com/opctl/opctl/cli/cmd/auth"
	"github.com/opctl/opctl/cli/cmd/node"
	"github.com/opctl/opctl/cli/cmd/op"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	cliColorer = clicolorer.New()
	cliOutput  = clioutput.New(cliColorer, os.Stderr, os.Stdout)
	nodeConfig local.NodeConfig
	noColor    bool
	version    string // set via -ldflags=-X=github.com/opctl/opctl/cli/cmd.version=$(version)
)

func printUsageSection(
	cmd *cobra.Command,
	heading string,
	lines []string,
) {
	// don't print empty sections
	if len(lines) == 0 {
		return
	}

	// don't print empty headings
	if heading != "" {
		fmt.Fprintf(
			cmd.OutOrStderr(),
			"%s:\n",
			heading,
		)
	}

	for _, line := range lines {
		prefix := ""
		// only indent sections with headings
		if heading != "" {
			prefix = "  "
		}

		fmt.Fprintf(
			cmd.OutOrStderr(),
			"%s%s\n",
			prefix,
			line,
		)
	}

	fmt.Fprint(
		cmd.OutOrStderr(),
		"\n",
	)
}

func NewRootCmd() (*cobra.Command, error) {
	cliParamSatisfier := cliparamsatisfier.New(
		cliOutput,
	)

	rootCmd := cobra.Command{
		// hide completion options until we make them useful; defaults aren't
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		DisableAutoGenTag: true,
		Use:               "opctl",
		Short:             "Opctl is a free and open source distributed operation control system",
		Version:           version,
		// we do our own error printing
		SilenceErrors: true,
	}

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if noColor {
			cliColorer.DisableColor()
		}
		// don't include usage in non flag/argument errors
		cmd.SilenceUsage = true
	}

	rootCmd.SetUsageFunc(
		func(cmd *cobra.Command) error {
			if cmd.Runnable() {
				printUsageSection(
					cmd,
					"Usage",
					[]string{
						cmd.UseLine(),
					},
				)

				if strings.TrimSpace(cmd.Long) != "" {
					printUsageSection(
						cmd,
						"Synopsis",
						strings.Split(cmd.Long, "\n"),
					)
				}
			}
			if cmd.HasAvailableSubCommands() {
				printUsageSection(
					cmd,
					"Usage",
					[]string{
						fmt.Sprintf(
							"%s [command]",
							cmd.CommandPath(),
						),
					},
				)
			}

			if cmd.HasAvailableLocalFlags() {
				printUsageSection(
					cmd,
					"Flags",
					strings.Split(cmd.LocalFlags().FlagUsages(), "\n"),
				)
			}
			if cmd.HasInheritedFlags() {
				printUsageSection(
					cmd,
					"Global Flags",
					strings.Split(cmd.InheritedFlags().FlagUsages(), "\n"),
				)
			}
			if strings.TrimSpace(cmd.Example) != "" {
				printUsageSection(
					cmd,
					"Examples",
					strings.Split(cmd.Example, "\n"),
				)
			}
			if len(cmd.Commands()) > 0 {
				lines := []string{}
				for _, subCmd := range cmd.Commands() {
					if subCmd.IsAvailableCommand() {
						lines = append(
							lines,
							fmt.Sprintf(
								"%-*s   %s",
								subCmd.NamePadding(),
								subCmd.Name(),
								subCmd.Short),
						)
					}
				}
				printUsageSection(
					cmd,
					"Available Commands",
					lines,
				)
			}
			if cmd.HasAvailableSubCommands() {
				printUsageSection(
					cmd,
					"",
					[]string{
						fmt.Sprintf(
							"Use \"%s [command] --help\" for more information about a command.",
							cmd.CommandPath(),
						),
					},
				)
			}

			return nil
		},
	)

	rootCmd.SetHelpTemplate(
		`{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`,
	)

	// we don't want a version subcommand; we have a --help flag
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// we want -v to return the version only
	rootCmd.SetVersionTemplate(`{{printf "%s" .Version}}
`)

	perUserAppDataPath, err := appdatapath.New().PerUser()
	if err != nil {
		return nil, err
	}

	// persistent flags
	rootCmd.
		PersistentFlags().
		StringVar(
			&nodeConfig.APIListenAddress,
			"api-listen-address",
			"127.0.0.1:42224",
			"IP:PORT on which the API server will listen",
		)

	rootCmd.
		PersistentFlags().
		StringVar(
			&nodeConfig.ContainerRuntime,
			"container-runtime",
			"docker",
			"Runtime for opctl containers. Can be 'docker', 'k8s', or 'qemu' (experimental)",
		)

	rootCmd.
		PersistentFlags().
		StringVar(
			&nodeConfig.DataDir,
			"data-dir",
			filepath.Join(perUserAppDataPath, "opctl"),
			"Path of dir used to store opctl data",
		)

	rootCmd.
		PersistentFlags().
		StringVar(
			&nodeConfig.DNSListenAddress,
			"dns-listen-address",
			"127.0.0.1:53",
			"IP:PORT on which the DNS server will listen",
		)

	rootCmd.
		PersistentFlags().
		BoolVar(
			&noColor,
			"no-color",
			false,
			"Disable output coloring",
		)

	// add commands
	rootCmd.AddCommand(
		auth.NewAuthCmd(
			&nodeConfig,
		),
	)
	rootCmd.AddCommand(
		newEventsCmd(
			cliOutput,
			&nodeConfig,
		),
	)
	rootCmd.AddCommand(
		newLsCmd(
			cliParamSatisfier,
			&nodeConfig,
		),
	)
	rootCmd.AddCommand(
		node.NewNodeCmd(
			cliColorer,
			&nodeConfig,
		),
	)
	rootCmd.AddCommand(
		op.NewOpCmd(
			cliParamSatisfier,
			&nodeConfig,
		),
	)
	rootCmd.AddCommand(
		newRunCmd(
			cliOutput,
			cliParamSatisfier,
		),
	)
	rootCmd.AddCommand(
		newSelfUpdateCmd(
			cliOutput,
			&nodeConfig,
		),
	)
	rootCmd.AddCommand(
		newUICmd(
			cliOutput,
			cliParamSatisfier,
			&nodeConfig,
		),
	)

	populateFlagsFromEnvVars(&rootCmd)

	return &rootCmd, nil
}

func populateFlagsFromEnvVars(
	cmd *cobra.Command,
) {
	envVarPrefix := strings.ReplaceAll(
		strings.ToUpper(
			cmd.CommandPath(),
		),
		" ",
		"_",
	)

	cmd.
		LocalFlags().
		VisitAll(
			func(f *pflag.Flag) {
				envVarName := strings.Join(
					[]string{
						envVarPrefix,
						strings.ToUpper(
							strings.ReplaceAll(f.Name, "-", "_"),
						),
					},
					"_",
				)

				envVarValue := os.Getenv(envVarName)

				// include env vars in usage docs
				f.Usage = fmt.Sprintf("%s (env $%s)", f.Usage, envVarName)

				// apply the env var value to the flag when the flag is not set and the env var is
				if !f.Changed && envVarValue != "" {
					cmd.Flags().Set(f.Name, envVarValue)
				}
			},
		)

	for _, childCmd := range cmd.Commands() {
		populateFlagsFromEnvVars(childCmd)
	}
}

func Execute() {
	rootCmd, err := NewRootCmd()
	if err != nil {
		cliOutput.Error(err.Error())
		os.Exit(1)
	}

	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		cliOutput.Error(err.Error())

		if re, ok := err.(clioutput.RunError); ok {
			os.Exit(re.ExitCode)
		} else {
			os.Exit(1)
		}
	}
}
