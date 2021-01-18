package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/appdataspec/sdk-golang/appdatapath"
	mow "github.com/jawher/mow.cli"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/datadir"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/cli/internal/updater"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/core"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime/docker"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime/k8s"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"github.com/skratchdot/open-golang/open"
	"os/signal"
	"text/tabwriter"
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

			resources := addCmd.StringArg("RESOURCES", "", "Resources this auth applies to in the form of a host or host/path.")
			username := addCmd.StringOpt("u username", "", "Username")
			password := addCmd.StringOpt("p password", "", "Password")

			addCmd.Action = func() {
				node, err := nodeProvider.CreateNodeIfNotExists()
				if nil != err {
					exitWith("", err)
				}

				exitWith(
					"",
					node.AddAuth(
						context.TODO(),
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
			node, err := nodeProvider.CreateNodeIfNotExists()
			if err != nil {
				exitWith("", err)
			}

			eventChannel, err := node.GetEventStream(
				context.TODO(),
				&model.GetEventStreamReq{},
			)
			if nil != err {
				exitWith("", err)
			}

			for {
				event, isEventChannelOpen := <-eventChannel
				if !isEventChannelOpen {
					exitWith("", errors.New("Connection to event stream lost"))
				}

				cliOutput.Event(&event)
			}
		}
	})

	cli.Command("ls", "List operations", func(lsCmd *mow.Cmd) {
		const dirRefArgName = "DIR_REF"
		lsCmd.Spec = fmt.Sprintf("[%v]", dirRefArgName)
		dirRef := lsCmd.StringArg(dirRefArgName, opspec.DotOpspecDirName, "Reference to dir ops will be listed from")

		lsCmd.Action = func() {
			node, err := nodeProvider.CreateNodeIfNotExists()
			if err != nil {
				exitWith("", err)
			}

			dataResolver := dataresolver.New(
				cliParamSatisfier,
				node,
			)

			_tabWriter := new(tabwriter.Writer)
			defer _tabWriter.Flush()
			_tabWriter.Init(os.Stdout, 0, 8, 1, '\t', 0)

			fmt.Fprintln(_tabWriter, "REF\tDESCRIPTION")

			dirHandle, err := dataResolver.Resolve(
				*dirRef,
				nil,
			)
			if nil != err {
				exitWith("", err)
			}

			opsByPath, err := opspec.List(
				context.TODO(),
				dirHandle,
			)
			if nil != err {
				exitWith("", err)
			}

			cwd, err := os.Getwd()
			if nil != err {
				exitWith("", err)
			}

			for path, op := range opsByPath {
				opRef := filepath.Join(dirHandle.Ref(), path)
				if filepath.IsAbs(opRef) {
					// make absolute paths relative
					relOpRef, err := filepath.Rel(cwd, opRef)
					if nil != err {
						exitWith("", err)
					}

					opRef = strings.TrimPrefix(relOpRef, ".opspec/")
				}

				fmt.Fprintf(_tabWriter, "%v\t%v", opRef, op.Description)
				fmt.Fprintln(_tabWriter)
			}
		}
	})

	cli.Command("node", "Manage nodes", func(nodeCmd *mow.Cmd) {
		nodeCmd.Command("create", "Creates a node", func(createCmd *mow.Cmd) {
			createCmd.Action = func() {

				dataDir, err := datadir.New(nodeCreateOpts.DataDir)
				if nil != err {
					exitWith("", err)
				}

				if err := dataDir.InitAndLock(); nil != err {
					exitWith("", err)
				}

				var containerRuntime containerruntime.ContainerRuntime
				if "k8s" == nodeCreateOpts.ContainerRuntime {
					containerRuntime, err = k8s.New()
					if nil != err {
						exitWith("", err)
					}
				} else {
					containerRuntime, err = docker.New()
					if nil != err {
						exitWith("", err)
					}
				}

				exitWith(
					"",
					newHTTPListener(
						core.New(
							containerRuntime,
							dataDir.Path(),
						),
					).
						listen(
							context.Background(),
							nodeCreateOpts.ListenAddress,
						),
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
				// install the whole pkg in case relative (intra pkg) refs exist
				opRefParts := strings.SplitN(*opRef, "#", 2)
				var pkgRef string
				if len(opRefParts) == 1 {
					pkgRef = opRefParts[0]
				} else {
					if verAndPathParts := strings.SplitN(opRefParts[1], "/", 2); len(verAndPathParts) != 1 {
						pkgRef = fmt.Sprintf("%s#%s", opRefParts[0], verAndPathParts[0])
					}
				}

				opDirHandle, err := dataResolver.Resolve(
					pkgRef,
					&model.Creds{
						Username: *username,
						Password: *password,
					},
				)
				if err != nil {
					exitWith("", err)
				}

				exitWith(
					"",
					opspec.Install(
						context.TODO(),
						filepath.Join(*path, pkgRef),
						opDirHandle,
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
				opDirHandle, err := dataResolver.Resolve(
					*opRef,
					nil,
				)
				if nil != err {
					exitWith("", err)
				}

				exitWith(
					fmt.Sprintf("%v is valid", opDirHandle.Ref()),
					opspec.Validate(context.TODO(), *opDirHandle.Path()),
				)
			}
		})
	})

	cli.Command("run", "Start and wait on an op", func(runCmd *mow.Cmd) {
		args := runCmd.StringsOpt("a", []string{}, "Explicitly pass args to op in format `-a NAME1=VALUE1 -a NAME2=VALUE2`")
		argFile := runCmd.StringOpt("arg-file", filepath.Join(opspec.DotOpspecDirName, "args.yml"), "Read in a file of args in yml format")
		opRef := runCmd.StringArg("OP_REF", "", "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)")

		runCmd.Action = func() {
			startTime := time.Now().UTC()
			ctx := context.TODO()

			node, err := nodeProvider.CreateNodeIfNotExists()
			if err != nil {
				exitWith("", err)
			}

			dataResolver := dataresolver.New(
				cliParamSatisfier,
				node,
			)

			opHandle, err := dataResolver.Resolve(
				*opRef,
				nil,
			)
			if nil != err {
				exitWith("", err)
			}

			opFileReader, err := opHandle.GetContent(
				ctx,
				opfile.FileName,
			)
			if nil != err {
				exitWith("", err)
			}

			opFileBytes, err := ioutil.ReadAll(opFileReader)
			if nil != err {
				exitWith("", err)
			}

			opFile, err := opfile.Unmarshal(
				opFileBytes,
			)
			if nil != err {
				exitWith("", err)
			}

			ymlFileInputSrc, err := cliParamSatisfier.NewYMLFileInputSrc(*argFile)
			if nil != err {
				exitWith("", fmt.Errorf("unable to load arg file at '%v'; error was: %v", *argFile, err.Error()))
			}

			argsMap, err := cliParamSatisfier.Satisfy(
				cliparamsatisfier.NewInputSourcer(
					cliParamSatisfier.NewSliceInputSrc(*args, "="),
					ymlFileInputSrc,
					cliParamSatisfier.NewEnvVarInputSrc(),
					cliParamSatisfier.NewParamDefaultInputSrc(opFile.Inputs),
					cliParamSatisfier.NewCliPromptInputSrc(opFile.Inputs),
				),
				opFile.Inputs,
			)
			if nil != err {
				exitWith("", err)
			}

			// init signal channels
			aSigIntWasReceivedAlready := false
			sigIntChannel := make(chan os.Signal, 1)
			defer close(sigIntChannel)
			signal.Notify(
				sigIntChannel,
				syscall.SIGINT,
			)

			sigTermChannel := make(chan os.Signal, 1)
			defer close(sigTermChannel)
			signal.Notify(
				sigTermChannel,
				syscall.SIGTERM,
			)

			// start op
			rootCallID, err := node.StartOp(
				ctx,
				model.StartOpReq{
					Args: argsMap,
					Op: model.StartOpReqOp{
						Ref: opHandle.Ref(),
					},
				},
			)
			if nil != err {
				exitWith("", err)
			}

			// start event loop
			eventChannel, err := node.GetEventStream(
				ctx,
				&model.GetEventStreamReq{
					Filter: model.EventFilter{
						Roots: []string{rootCallID},
						Since: &startTime,
					},
				},
			)
			if nil != err {
				exitWith("", err)
			}

			for {
				select {

				case <-sigIntChannel:
					if !aSigIntWasReceivedAlready {
						cliOutput.Warning("Gracefully stopping... (signal Control-C again to force)")
						aSigIntWasReceivedAlready = true

						node.KillOp(
							ctx,
							model.KillOpReq{
								OpID:       rootCallID,
								RootCallID: rootCallID,
							},
						)
					} else {
						exitWith("", &RunError{
							ExitCode: 130,
							message:  "Terminated by Control-C",
						})
					}

				case <-sigTermChannel:
					cliOutput.Warning("Gracefully stopping...")

					exitWith(
						"",
						node.KillOp(
							ctx,
							model.KillOpReq{
								OpID:       rootCallID,
								RootCallID: rootCallID,
							},
						),
					)
				case event, isEventChannelOpen := <-eventChannel:
					if !isEventChannelOpen {
						exitWith("", errors.New("Event channel closed unexpectedly"))
					}

					cliOutput.Event(&event)

					if nil != event.CallEnded {
						if event.CallEnded.Call.ID == rootCallID {
							switch event.CallEnded.Outcome {
							case model.OpOutcomeSucceeded:
								exitWith("", nil)
							case model.OpOutcomeKilled:
								exitWith("", &RunError{ExitCode: 137})
							default:
								exitWith("", &RunError{ExitCode: 1})
							}
						}
					}
				}
			}
		}
	})

	cli.Command("self-update", "Update opctl", func(selfUpdateCmd *mow.Cmd) {
		channel := selfUpdateCmd.StringOpt("c channel", "stable", "Release channel to update from (either `stable`, `alpha`, or `beta`)")
		selfUpdateCmd.Action = func() {
			updater := updater.New()
			update, err := updater.GetUpdateIfExists(*channel)
			if nil != err {
				exitWith("", err)
			} else if nil == update {
				exitWith("No update available, already at the latest version!", nil)
			}

			err = updater.ApplyUpdate(update)
			if nil != err {
				exitWith("", err)
			}

			// kill local node to ensure outdated version not left running
			// @TODO start node maintaining previous user
			err = nodeProvider.KillNodeIfExists("")
			if nil != err {
				err = fmt.Errorf("Unable to kill running node; run `node kill` to complete the update. Error was: %v", err)
			}
			exitWith(
				fmt.Sprintf("Updated to new version: %s!", update.Version),
				err,
			)
		}
	})

	cli.Command("ui", "Open the opctl web UI and mount a reference.", func(uiCmd *mow.Cmd) {
		const mountRefArgName = "MOUNT_REF"
		uiCmd.Spec = fmt.Sprintf("[%v]", mountRefArgName)
		mountRefArg := uiCmd.StringArg(mountRefArgName, ".", "Reference to mount (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)")

		uiCmd.Action = func() {
			var resolvedMount string
			var err error
			if strings.HasPrefix(*mountRefArg, ".") {
				// treat dot paths as regular rel paths
				resolvedMount, err = filepath.Abs(*mountRefArg)
				if nil != err {
					exitWith("", err)
				}
			} else {
				node, err := nodeProvider.CreateNodeIfNotExists()
				if err != nil {
					exitWith("", err)
				}

				dataResolver := dataresolver.New(
					cliParamSatisfier,
					node,
				)

				// otherwise use same resolution as run
				mountHandle, err := dataResolver.Resolve(
					*mountRefArg,
					nil,
				)
				if nil != err {
					exitWith("", err)
				}

				resolvedMount = mountHandle.Ref()
			}

			exitWith(
				"Opctl web UI opened!",
				open.Run(
					fmt.Sprintf("http://localhost:42224?mount=%s", url.QueryEscape(resolvedMount)),
				),
			)
		}
	})

	return cli
}
