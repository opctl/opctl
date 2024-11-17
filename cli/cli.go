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
	"github.com/opctl/opctl/cli/internal/oppath"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"golang.org/x/term"
)

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

			if re, ok := err.(RunError); ok {
				mow.Exit(re.ExitCode)
			} else {
				mow.Exit(1)
			}
		}
	}

	perUserAppDataPath, err := appdatapath.New().PerUser()
	if err != nil {
		exitWith("", err)
	}

	containerRuntime := cli.String(
		mow.StringOpt{
			Desc:   "Runtime for opctl containers. Can be 'docker', 'k8s', or 'qemu' (experimental)",
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

	apiListenAddress := cli.String(
		mow.StringOpt{
			Desc:   "IP:PORT on which the API server will listen",
			EnvVar: "OPCTL_API_LISTEN_ADDRESS",
			Name:   "api-listen-address",
			Value:  "127.0.0.1:42224",
		},
	)

	dnsListenAddress := cli.String(
		mow.StringOpt{
			Desc:   "IP:PORT on which the DNS server will listen",
			EnvVar: "OPCTL_DNS_LISTEN_ADDRESS",
			Name:   "dns-listen-address",
			Value:  "127.0.0.1:53",
		},
	)

	cliParamSatisfier := cliparamsatisfier.New(cliOutput)

	noColor := cli.Bool(mow.BoolOpt{
		Name:   "nc no-color",
		Value:  false,
		Desc:   "Disable output coloring",
		EnvVar: "OPCTL_NO_COLOR",
	})

	cli.Before = func() {
		if *noColor {
			cliOutput.DisableColor()
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	cli.After = func() {
		cancel()
	}

	cli.Command("auth", "Manage default auth for OCI image and git repositories", func(authCmd *mow.Cmd) {
		authCmd.Command("add", "Add default auth for OCI image and git repositories", func(addCmd *mow.Cmd) {
			addCmd.Spec = "RESOURCES [ -u=<username> ] [ -p=<password> ]"

			resources := addCmd.String(mow.StringArg{
				Name:   "RESOURCES",
				Value:  "",
				Desc:   "Resources this auth applies to in the form of a ref prefix (e.g. docker.io, github.com/some-org, etc.)",
				EnvVar: "OPCTL_AUTH_RESOURCES",
			})
			username := addCmd.String(mow.StringOpt{
				Name:   "u username",
				Value:  "",
				Desc:   "Username",
				EnvVar: "OPCTL_AUTH_USERNAME",
			})
			password := addCmd.String(mow.StringOpt{
				Name:   "p password",
				Value:  "",
				Desc:   "Password",
				EnvVar: "OPCTL_AUTH_PASSWORD",
			})

			addCmd.Action = func() {
				exitWith(
					"",
					auth(
						ctx,
						local.NodeConfig{
							APIListenAddress: *apiListenAddress,
							ContainerRuntime: *containerRuntime,
							DataDir:          *dataDir,
							DNSListenAddress: *dnsListenAddress,
						},
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
					ctx,
					cliOutput,
					local.NodeConfig{
						APIListenAddress: *apiListenAddress,
						ContainerRuntime: *containerRuntime,
						DataDir:          *dataDir,
						DNSListenAddress: *dnsListenAddress,
					},
				),
			)
		}
	})

	cli.Command("ls", "List operations", func(lsCmd *mow.Cmd) {
		const dirRefArgName = "DIR_REF"
		lsCmd.Spec = fmt.Sprintf("[%v]", dirRefArgName)

		dirRef := lsCmd.String(mow.StringArg{
			Name:   dirRefArgName,
			Value:  model.DotOpspecDirName,
			Desc:   "Reference to dir ops will be listed from",
			EnvVar: "OPCTL_LS_DIR_REF",
		})

		lsCmd.Action = func() {
			exitWith(
				"",
				ls(
					ctx,
					cliParamSatisfier,
					local.NodeConfig{
						APIListenAddress: *apiListenAddress,
						ContainerRuntime: *containerRuntime,
						DataDir:          *dataDir,
						DNSListenAddress: *dnsListenAddress,
					},
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
					nodeCreate(
						ctx,
						local.NodeConfig{
							APIListenAddress: *apiListenAddress,
							ContainerRuntime: *containerRuntime,
							DataDir:          *dataDir,
							DNSListenAddress: *dnsListenAddress,
						},
					),
				)
			}
		})

		nodeCmd.Command("delete", "Deletes a node. This is destructive! all node data including auth, caches, and operation state will be permanently removed.", func(deleteCmd *mow.Cmd) {
			deleteCmd.Action = func() {
				exitWith(
					"",
					nodeDelete(
						ctx,
						local.NodeConfig{
							APIListenAddress: *apiListenAddress,
							ContainerRuntime: *containerRuntime,
							DataDir:          *dataDir,
							DNSListenAddress: *dnsListenAddress,
						},
					),
				)
			}
		})

		nodeCmd.Command("kill", "Kills a node and any running operations. This is non destructive. All node data including auth, caches, and operation state will be retained.", func(killCmd *mow.Cmd) {
			killCmd.Action = func() {
				exitWith(
					"",
					nodeKill(
						ctx,
						local.NodeConfig{
							APIListenAddress: *apiListenAddress,
							ContainerRuntime: *containerRuntime,
							DataDir:          *dataDir,
							DNSListenAddress: *dnsListenAddress,
						},
					),
				)
			}
		})
	})

	cli.Command("op", "Manage ops", func(opCmd *mow.Cmd) {
		np, err := local.New(
			local.NodeConfig{
				APIListenAddress: *apiListenAddress,
				ContainerRuntime: *containerRuntime,
				DataDir:          *dataDir,
				DNSListenAddress: *dnsListenAddress,
			},
		)
		if err != nil {
			exitWith("", err)
		}

		node, err := np.CreateNodeIfNotExists(ctx)
		if err != nil {
			exitWith("", err)
		}

		dataResolver := dataresolver.New(
			cliParamSatisfier,
			node,
		)
		cwd, err := os.Getwd()
		if err != nil {
			exitWith("", err)
		}

		opPath, err := oppath.Get(
			ctx,
			cwd,
			dataResolver,
			node,
		)
		if err != nil {
			exitWith("", err)
		}

		opCmd.Command("create", "Create an op", func(createCmd *mow.Cmd) {
			path := createCmd.String(mow.StringOpt{
				Name:   "path",
				Value:  model.DotOpspecDirName,
				Desc:   "Path the op will be created at",
				EnvVar: "OPCTL_CREATE_PATH",
			})
			description := createCmd.String(mow.StringOpt{
				Name:   "d description",
				Value:  "",
				Desc:   "Op description",
				EnvVar: "OPCTL_CREATE_DESCRIPTION",
			})
			name := createCmd.String(mow.StringArg{
				Name:   "NAME",
				Value:  "",
				Desc:   "Op name",
				EnvVar: "OPCTL_CREATE_NAME",
			})

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
			path := installCmd.String(mow.StringOpt{
				Name:   "path",
				Value:  model.DotOpspecDirName,
				Desc:   "Path the op will be installed at",
				EnvVar: "OPCTL_INSTALL_PATH",
			})
			opRef := installCmd.String(mow.StringArg{
				Name:   "OP_REF",
				Value:  "",
				Desc:   "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)",
				EnvVar: "OPCTL_INSTALL_OP_REF",
			})
			username := installCmd.String(mow.StringOpt{
				Name:   "u username",
				Value:  "",
				Desc:   "Username used to auth w/ the pkg source",
				EnvVar: "OPCTL_INSTALL_USERNAME",
			})
			password := installCmd.String(mow.StringOpt{
				Name:   "p password",
				Value:  "",
				Desc:   "Password used to auth w/ the pkg source",
				EnvVar: "OPCTL_INSTALL_PASSWORD",
			})

			installCmd.Action = func() {
				exitWith(
					"",
					opInstall(
						ctx,
						dataResolver,
						opPath,
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
			opID := killCmd.String(mow.StringArg{
				Name:   "OP_ID",
				Value:  "",
				Desc:   "Id of the op to kill",
				EnvVar: "OPCTL_KILL_OP_ID",
			})

			killCmd.Action = func() {
				exitWith(
					"",
					node.KillOp(
						ctx,
						model.KillOpReq{
							OpID:       *opID,
							RootCallID: *opID,
						},
					),
				)
			}
		})

		opCmd.Command("validate", "Validate an op", func(validateCmd *mow.Cmd) {
			opRef := validateCmd.String(mow.StringArg{
				Name:   "OP_REF",
				Value:  "",
				Desc:   "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)",
				EnvVar: "OPCTL_VALIDATE_OP_REF",
			})

			validateCmd.Action = func() {
				exitWith(
					fmt.Sprintf("%v is valid", *opRef),
					opValidate(
						ctx,
						dataResolver,
						opPath,
						*opRef,
					),
				)
			}
		})
	})

	cli.Command("run", "Start and wait on an op", func(runCmd *mow.Cmd) {
		args := runCmd.Strings(mow.StringsOpt{
			Name:   "a",
			Desc:   "Explicitly pass args to op in format `-a NAME1=VALUE1 -a NAME2=VALUE2`",
			Value:  []string{},
			EnvVar: "OPCTL_RUN_ARGS",
		})
		argFile := runCmd.String(mow.StringOpt{
			Name:   "arg-file",
			Desc:   "Read in a file of args in yml format",
			Value:  filepath.Join(model.DotOpspecDirName, "args.yml"),
			EnvVar: "OPCTL_RUN_ARG_FILE",
		})
		noProgress := runCmd.Bool(mow.BoolOpt{
			Name:   "no-progress",
			Desc:   "Disable live call graph for the op",
			Value:  !term.IsTerminal(int(os.Stdout.Fd())),
			EnvVar: "OPCTL_RUN_NO_PROGRESS",
		})
		opRef := runCmd.String(mow.StringArg{
			Name:   "OP_REF",
			Desc:   "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)",
			Value:  "",
			EnvVar: "OPCTL_RUN_OP_REF",
		})

		runCmd.Action = func() {
			exitWith(
				"",
				run(
					ctx,
					cliOutput,
					cliParamSatisfier,
					local.NodeConfig{
						APIListenAddress: *apiListenAddress,
						ContainerRuntime: *containerRuntime,
						DataDir:          *dataDir,
						DNSListenAddress: *dnsListenAddress,
					},
					*args,
					*argFile,
					*opRef,
					*noProgress,
				),
			)
		}
	})

	cli.Command("self-update", "Update opctl", func(selfUpdateCmd *mow.Cmd) {
		selfUpdateCmd.Action = func() {
			exitWith(
				selfUpdate(
					ctx,
					local.NodeConfig{
						APIListenAddress: *apiListenAddress,
						ContainerRuntime: *containerRuntime,
						DataDir:          *dataDir,
						DNSListenAddress: *dnsListenAddress,
					},
				),
			)
		}
	})

	cli.Command("ui", "Open the opctl web UI and mount a reference.", func(uiCmd *mow.Cmd) {
		const mountRefArgName = "MOUNT_REF"
		uiCmd.Spec = fmt.Sprintf("[%v]", mountRefArgName)
		mountRefArg := uiCmd.String(mow.StringArg{
			Name:   mountRefArgName,
			Desc:   "Reference to mount (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)",
			Value:  ".",
			EnvVar: "OPCTL_UI_MOUNT_REF",
		})

		uiCmd.Action = func() {
			exitWith(
				"Opctl web UI opened!",
				ui(
					ctx,
					cliParamSatisfier,
					local.NodeConfig{
						APIListenAddress: *apiListenAddress,
						ContainerRuntime: *containerRuntime,
						DataDir:          *dataDir,
						DNSListenAddress: *dnsListenAddress,
					},
					*apiListenAddress,
					*mountRefArg,
				),
			)
		}
	})

	return cli
}
