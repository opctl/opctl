package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/appdataspec/sdk-golang/appdatapath"
	mow "github.com/jawher/mow.cli"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	mowclitypes "github.com/opctl/opctl/cli/internal/mowclitypes"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/cli/internal/telemetry"
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

	containerRuntime := cli.String(
		mow.StringOpt{
			Desc:   "Runtime for opctl containers. Can be 'docker', 'k8s', or 'qemu' (experimental)",
			EnvVar: "OPCTL_CONTAINER_RUNTIME",
			Name:   "container-runtime",
			Value:  "docker",
		},
	)

	perUserAppDataPath, perUserAppDataPathErr := appdatapath.New().PerUser()
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

	ctx, cancel := context.WithCancel(context.Background())
	t, err := telemetry.Setup(ctx, *dataDir)
	if err != nil {
		cliOutput.Error(err.Error())
	}

	exitWith := func(successMessage string, err error) {
		if t != nil {
			t.Shutdown(ctx)
		}

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

	instrumentExitWithAction := func(name string, action func() (successMessage string, err error)) func() {
		return func() {
			start := time.Now()
			telemetry.InstrumentCommandUsage(ctx, name)
			successMessage, err := action()
			if err != nil {
				telemetry.InstrumentErrorMetrics(ctx, err)
			}
			telemetry.InstrumentExecutionTime(ctx, name, start)
			exitWith(successMessage, err)
		}
	}

	if perUserAppDataPathErr != nil {
		exitWith("", perUserAppDataPathErr)
	}

	cliParamSatisfier := cliparamsatisfier.New(cliOutput)

	noColor := cli.BoolOpt("nc no-color", false, "Disable output coloring")
	cli.Before = func() {
		if *noColor {
			cliOutput.DisableColor()
		}
	}

	cli.After = func() {
		cancel()
	}

	cli.Command("auth", "Manage auth for OCI image registries", func(authCmd *mow.Cmd) {
		authCmd.Command("add", "Add auth for an OCI image registry", func(addCmd *mow.Cmd) {
			addCmd.Spec = "RESOURCES [ -u=<username> ] [ -p=<password> ]"

			resources := addCmd.StringArg("RESOURCES", "", "Resources this auth applies to in the form of a host or host/path (e.g. docker.io)")
			username := addCmd.StringOpt("u username", "", "Username")
			password := addCmd.StringOpt("p password", "", "Password")

			addCmd.Action = instrumentExitWithAction("auth.add", func() (string, error) {
				return "", auth(
					ctx,
					local.NodeConfig{
						ContainerRuntime: *containerRuntime,
						DataDir:          *dataDir,
						ListenAddress:    *listenAddress,
					},
					model.AddAuthReq{
						Resources: *resources,
						Creds: model.Creds{
							Username: *username,
							Password: *password,
						},
					},
				)
			})
		})
	})

	cli.Command("events", "Stream events", func(eventsCmd *mow.Cmd) {
		eventsCmd.Action = instrumentExitWithAction("events", func() (string, error) {
			return "", events(
				ctx,
				cliOutput,
				local.NodeConfig{
					ContainerRuntime: *containerRuntime,
					DataDir:          *dataDir,
					ListenAddress:    *listenAddress,
				},
			)
		})
	})

	cli.Command("ls", "List operations", func(lsCmd *mow.Cmd) {
		const dirRefArgName = "DIR_REF"
		lsCmd.Spec = fmt.Sprintf("[%v]", dirRefArgName)
		dirRef := lsCmd.StringArg(dirRefArgName, opspec.DotOpspecDirName, "Reference to dir ops will be listed from")

		lsCmd.Action = instrumentExitWithAction("ls", func() (string, error) {
			return "", ls(
				ctx,
				cliParamSatisfier,
				local.NodeConfig{
					ContainerRuntime: *containerRuntime,
					DataDir:          *dataDir,
					ListenAddress:    *listenAddress,
				},
				*dirRef,
			)
		})
	})

	cli.Command("node", "Manage nodes", func(nodeCmd *mow.Cmd) {
		nodeCmd.Command("create", "Creates a node", func(createCmd *mow.Cmd) {
			createCmd.Action = instrumentExitWithAction("node.create", func() (string, error) {
				return "", nodeCreate(
					ctx,
					local.NodeConfig{
						ContainerRuntime: *containerRuntime,
						DataDir:          *dataDir,
						ListenAddress:    *listenAddress,
					},
				)
			})
		})

		nodeCmd.Command("delete", "Deletes a node. This is destructive! all node data including auth, caches, and operation state will be permanently removed.", func(deleteCmd *mow.Cmd) {
			deleteCmd.Action = instrumentExitWithAction("node.delete", func() (string, error) {
				return "", nodeDelete(
					ctx,
					local.NodeConfig{
						ContainerRuntime: *containerRuntime,
						DataDir:          *dataDir,
						ListenAddress:    *listenAddress,
					},
				)
			})
		})

		nodeCmd.Command("kill", "Kills a node and any running operations. This is non destructive. All node data including auth, caches, and operation state will be retained.", func(killCmd *mow.Cmd) {
			killCmd.Action = instrumentExitWithAction("node.kill", func() (string, error) {
				return "", nodeKill(
					ctx,
					local.NodeConfig{
						ContainerRuntime: *containerRuntime,
						DataDir:          *dataDir,
						ListenAddress:    *listenAddress,
					},
				)
			})
		})
	})

	cli.Command("op", "Manage ops", func(opCmd *mow.Cmd) {
		np, err := local.New(
			local.NodeConfig{
				ContainerRuntime: *containerRuntime,
				DataDir:          *dataDir,
				ListenAddress:    *listenAddress,
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

		opCmd.Command("create", "Create an op", func(createCmd *mow.Cmd) {
			path := createCmd.StringOpt("path", opspec.DotOpspecDirName, "Path the op will be created at")
			description := createCmd.StringOpt("d description", "", "Op description")
			name := createCmd.StringArg("NAME", "", "Op name")

			createCmd.Action = instrumentExitWithAction("op.create", func() (string, error) {
				return "", opspec.Create(
					filepath.Join(*path, *name),
					*name,
					*description,
				)
			})
		})

		opCmd.Command("install", "Install an op", func(installCmd *mow.Cmd) {
			path := installCmd.StringOpt("path", opspec.DotOpspecDirName, "Path the op will be installed at")
			opRef := installCmd.StringArg("OP_REF", "", "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)")
			username := installCmd.StringOpt("u username", "", "Username used to auth w/ the pkg source")
			password := installCmd.StringOpt("p password", "", "Password used to auth w/ the pkg source")

			installCmd.Action = instrumentExitWithAction("op.install", func() (string, error) {
				return "", opInstall(
					ctx,
					dataResolver,
					*opRef,
					*path,
					&model.Creds{
						Username: *username,
						Password: *password,
					},
				)
			})
		})

		opCmd.Command("kill", "Kill an op", func(killCmd *mow.Cmd) {
			opID := killCmd.StringArg("OP_ID", "", "Id of the op to kill")

			killCmd.Action = instrumentExitWithAction("op.kill", func() (string, error) {
				return "", node.KillOp(
					ctx,
					model.KillOpReq{
						OpID:       *opID,
						RootCallID: *opID,
					},
				)
			})
		})

		opCmd.Command("validate", "Validate an op", func(validateCmd *mow.Cmd) {
			opRef := validateCmd.StringArg("OP_REF", "", "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)")

			validateCmd.Action = instrumentExitWithAction("op.validate", func() (string, error) {
				return fmt.Sprintf("%v is valid", *opRef), opValidate(
					ctx,
					dataResolver,
					*opRef,
				)
			})
		})
	})

	cli.Command("run", "Start and wait on an op", func(runCmd *mow.Cmd) {
		args := runCmd.StringsOpt("a", []string{}, "Explicitly pass args to op in format `-a NAME1=VALUE1 -a NAME2=VALUE2`")
		argFile := runCmd.StringOpt("arg-file", filepath.Join(opspec.DotOpspecDirName, "args.yml"), "Read in a file of args in yml format")
		noProgress := runCmd.BoolOpt("no-progress", !term.IsTerminal(int(os.Stdout.Fd())), "Disable live call graph for the op")
		opRef := runCmd.StringArg("OP_REF", "", "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)")

		runCmd.Action = instrumentExitWithAction("run", func() (string, error) {
			return "", run(
				ctx,
				cliOutput,
				cliParamSatisfier,
				local.NodeConfig{
					ContainerRuntime: *containerRuntime,
					DataDir:          *dataDir,
					ListenAddress:    *listenAddress,
				},
				*args,
				*argFile,
				*opRef,
				*noProgress,
			)
		})
	})

	cli.Command("self-update", "Update opctl", func(selfUpdateCmd *mow.Cmd) {
		selfUpdateCmd.Action = instrumentExitWithAction("self-update", func() (string, error) {
			return selfUpdate(
				local.NodeConfig{
					ContainerRuntime: *containerRuntime,
					DataDir:          *dataDir,
					ListenAddress:    *listenAddress,
				},
			)
		})
	})

	cli.Command("ui", "Open the opctl web UI and mount a reference.", func(uiCmd *mow.Cmd) {
		const mountRefArgName = "MOUNT_REF"
		uiCmd.Spec = fmt.Sprintf("[%v]", mountRefArgName)
		mountRefArg := uiCmd.StringArg(mountRefArgName, ".", "Reference to mount (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)")

		uiCmd.Action = instrumentExitWithAction("ui", func() (string, error) {
			return "Opctl web UI opened!", ui(
				ctx,
				cliParamSatisfier,
				local.NodeConfig{
					ContainerRuntime: *containerRuntime,
					DataDir:          *dataDir,
					ListenAddress:    *listenAddress,
				},
				*listenAddress,
				*mountRefArg,
			)
		})
	})

	cli.Command("telemetry", "Manage telemetry configuration", func(telemetryCmd *mow.Cmd) {
		telemetryCmd.Action = instrumentExitWithAction("telemetry", func() (string, error) {
			return telemetryGet(ctx, cliOutput, *dataDir)
		})

		telemetryCmd.Command("set", "Set telemetry configuration", func(setCmd *mow.Cmd) {
			setCmd.Spec = "[(-c=<collector> -e=<endpoint>)] [-i | -t | -m | -l]..."

			config, err := telemetry.ReadConfigFile(*dataDir)
			if err != nil {
				exitWith("", err)
				return
			}

			collector := mowclitypes.NewEnum(
				config.Collector,
				mowclitypes.EnumValue{
					Value: telemetry.PromPushGatewayCollector,
					Desc:  "Prometheus Pushgateway",
				},
			)
			setCmd.Var(mow.VarOpt{
				Name:   "c collector",
				Desc:   "Telemetry collector",
				Value:  collector,
				EnvVar: telemetry.EnvVarOpctlTelemetryCollector,
			})
			endpoint := setCmd.String(mow.StringOpt{
				Name:   "e endpoint",
				Desc:   "Telemetry HTTP endpoint",
				Value:  config.Endpoint,
				EnvVar: telemetry.EnvVarOpctlTelemetryEndpoint,
			})
			pushInterval := setCmd.IntOpt("i push-interval", int(config.PushInterval/time.Second), "Interval (in seconds to push telemetry data")
			tracesEnabled := setCmd.BoolOpt("t traces-enabled", config.TracesEnabled, "Enable traces")
			metricsEnabled := setCmd.BoolOpt("m metrics-enabled", config.MetricsEnabled, "Enable metrics")
			logsEnabled := setCmd.BoolOpt("l logs-enabled", config.LogsEnabled, "Enable logs")

			setCmd.Action = instrumentExitWithAction("telemetry.set", func() (string, error) {
				return telemetrySet(
					ctx,
					config,
					collector.Value,
					*endpoint,
					*dataDir,
					*pushInterval,
					*tracesEnabled,
					*metricsEnabled,
					*logsEnabled,
				)
			})
		})

		telemetryCmd.Command("push", "Manually push collected telemetry data", func(pushCmd *mow.Cmd) {
			pushCmd.Action = instrumentExitWithAction("telemetry.push", func() (string, error) {
				return "Pushed telemetry data!", telemetryPush(t, *dataDir)
			})
		})
	})

	return cli
}
