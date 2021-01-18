package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/appdataspec/sdk-golang/appdatapath"
	mow "github.com/jawher/mow.cli"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
)

var testModeEnvVar = "OPCTL_TEST_MODE"

type cli interface {
	Run(args []string) error
}

func newCli(
	cliOutput clioutput.CliOutput,
) cli {

	cli := mow.App(
		"opctl",
		"Opctl is a free and open source distributed operation control system.",
	)
	cli.Version("v version", version)

	exitWith := func(successMessage string, err error) {
		if err == nil {
			if successMessage != "" {
				cliOutput.Success(successMessage)
			}
			mow.Exit(0)
		} else {
			cliOutput.Error(err.Error())

			// @TODO find a better way to support tests & remove this
			// currently mow.Exit(non-zero) blows up test harness
			if _, isTestMode := os.LookupEnv(testModeEnvVar); isTestMode {
				os.Exit(0)
			}

			if re, ok := err.(*RunError); ok {
				mow.Exit(re.ExitCode)
			} else {
				mow.Exit(1)
			}
		}
	}

	perUserAppDataPath, err := appdatapath.New().PerUser()
	if nil != err {
		exitWith("", err)
	}

	containerRuntime := cli.String(
		mow.StringOpt{
			Desc:   "Runtime for opctl containers",
			EnvVar: "OPCTL_CONTAINER_RUNTIME",
			Name:   "container-runtime",
			Value:  "docker",
		},
	)

	dataDir := cli.String(
		mow.StringOpt{
			Desc:   "Path of dir used to store opctl data",
			EnvVar: "OPCTL_DATA_DIR",
			Name:   "data-dir",
			Value:  filepath.Join(perUserAppDataPath, "opctl"),
		},
	)

	listenAddress := cli.String(
		mow.StringOpt{
			Desc:   "HOST:PORT on which the node will listen",
			EnvVar: "OPCTL_LISTEN_ADDRESS",
			Name:   "listen-address",
			Value:  "127.0.0.1:42224",
		},
	)

	nodeCreateOpts := local.NodeCreateOpts{
		ContainerRuntime: *containerRuntime,
		DataDir:          *dataDir,
		ListenAddress:    *listenAddress,
	}

	nodeProvider := local.New(
		nodeCreateOpts,
	)

	cliParamSatisfier := cliparamsatisfier.New(cliOutput)

	noColor := cli.BoolOpt("nc no-color", false, "Disable output coloring")

	cli.Before = func() {
		if *noColor {
			cliOutput.DisableColor()
		}
	}

	cli.Command("auth", "Manage auth for OCI image registries", func(authCmd *mow.Cmd) {
		authCmd.Command("add", "Add auth for an OCI image registry", func(addCmd *mow.Cmd) {
			addCmd.Spec = "RESOURCES [ -u=<username> ] [ -p=<password> ]"

			resources := addCmd.StringArg("RESOURCES", "", "Resources this auth applies to in the form of a host or host/path (e.g. docker.io)")
			username := addCmd.StringOpt("u username", "", "Username")
			password := addCmd.StringOpt("p password", "", "Password")

			addCmd.Action = func() {
				exitWith(
					"",
					auth(
						context.TODO(),
						nodeProvider,
						model.AddAuthReq{
							Resources: *resources,
							Creds: model.Creds{
								Username: *username,
								Password: *password,
							},
						},
					),
				)
			}
		})
	})

	cli.Command("events", "Stream events", func(eventsCmd *mow.Cmd) {
		eventsCmd.Action = func() {
			exitWith(
				"",
				events(
					context.TODO(),
					cliOutput,
					nodeProvider,
				),
			)
		}
	})

	cli.Command("ls", "List operations", func(lsCmd *mow.Cmd) {
		const dirRefArgName = "DIR_REF"
		lsCmd.Spec = fmt.Sprintf("[%v]", dirRefArgName)
		dirRef := lsCmd.StringArg(dirRefArgName, opspec.DotOpspecDirName, "Reference to dir ops will be listed from")

		lsCmd.Action = func() {
			exitWith(
				"",
				ls(
					context.TODO(),
					cliParamSatisfier,
					nodeProvider,
					*dirRef,
				),
			)
		}
	})

	cli.Command("node", "Manage nodes", func(nodeCmd *mow.Cmd) {
		nodeCmd.Command("create", "Creates a node", func(createCmd *mow.Cmd) {
			createCmd.Action = func() {
				exitWith(
					"",
					node(nodeCreateOpts),
				)
			}
		})

		nodeCmd.Command("kill", "Kills a node", func(killCmd *mow.Cmd) {
			killCmd.Action = func() {
				exitWith("", nodeProvider.KillNodeIfExists(""))
			}
		})
	})

	cli.Command("op", "Manage ops", func(opCmd *mow.Cmd) {
		node, err := nodeProvider.CreateNodeIfNotExists()
		if err != nil {
			exitWith("", err)
		}

		dataResolver := dataresolver.New(
			cliParamSatisfier,
			node,
		)

		opCmd.Command("create", "Create an op", func(createCmd *mow.Cmd) {
			path := createCmd.StringOpt("path", opspec.DotOpspecDirName, "Path the op will be created at")
			description := createCmd.StringOpt("d description", "", "Op description")
			name := createCmd.StringArg("NAME", "", "Op name")

			createCmd.Action = func() {
				exitWith(
					"",
					opspec.Create(
						filepath.Join(*path, *name),
						*name,
						*description,
					),
				)
			}
		})

		opCmd.Command("install", "Install an op", func(installCmd *mow.Cmd) {
			path := installCmd.StringOpt("path", opspec.DotOpspecDirName, "Path the op will be installed at")
			opRef := installCmd.StringArg("OP_REF", "", "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)")
			username := installCmd.StringOpt("u username", "", "Username used to auth w/ the pkg source")
			password := installCmd.StringOpt("p password", "", "Password used to auth w/ the pkg source")

			installCmd.Action = func() {
				exitWith(
					"",
					opInstall(
						context.TODO(),
						dataResolver,
						*opRef,
						*path,
						&model.Creds{
							Username: *username,
							Password: *password,
						},
					),
				)
			}
		})

		opCmd.Command("kill", "Kill an op", func(killCmd *mow.Cmd) {
			opID := killCmd.StringArg("OP_ID", "", "Id of the op to kill")

			killCmd.Action = func() {
				exitWith(
					"",
					node.KillOp(
						context.TODO(),
						model.KillOpReq{
							OpID:       *opID,
							RootCallID: *opID,
						},
					),
				)
			}
		})

		opCmd.Command("validate", "Validate an op", func(validateCmd *mow.Cmd) {
			opRef := validateCmd.StringArg("OP_REF", "", "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)")

			validateCmd.Action = func() {
				exitWith(
					fmt.Sprintf("%v is valid", *opRef),
					opValidate(
						context.TODO(),
						dataResolver,
						*opRef,
					),
				)
			}
		})
	})

	cli.Command("run", "Start and wait on an op", func(runCmd *mow.Cmd) {
		args := runCmd.StringsOpt("a", []string{}, "Explicitly pass args to op in format `-a NAME1=VALUE1 -a NAME2=VALUE2`")
		argFile := runCmd.StringOpt("arg-file", filepath.Join(opspec.DotOpspecDirName, "args.yml"), "Read in a file of args in yml format")
		opRef := runCmd.StringArg("OP_REF", "", "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)")

		runCmd.Action = func() {
			exitWith(
				"",
				run(
					context.TODO(),
					cliOutput,
					cliParamSatisfier,
					nodeProvider,
					*args,
					*argFile,
					*opRef,
				),
			)
		}
	})

	cli.Command("self-update", "Update opctl", func(selfUpdateCmd *mow.Cmd) {
		channel := selfUpdateCmd.StringOpt("c channel", "stable", "Release channel to update from (either `stable`, `alpha`, or `beta`)")
		selfUpdateCmd.Action = func() {
			exitWith(
				selfUpdate(
					nodeProvider,
					*channel,
				),
			)
		}
	})

	cli.Command("ui", "Open the opctl web UI and mount a reference.", func(uiCmd *mow.Cmd) {
		const mountRefArgName = "MOUNT_REF"
		uiCmd.Spec = fmt.Sprintf("[%v]", mountRefArgName)
		mountRefArg := uiCmd.StringArg(mountRefArgName, ".", "Reference to mount (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)")

		uiCmd.Action = func() {
			exitWith(
				"Opctl web UI opened!",
				ui(
					cliParamSatisfier,
					nodeProvider,
					*mountRefArg,
				),
			)
		}
	})

	return cli
}
