package main

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/appdataspec/sdk-golang/appdatapath"
	mow "github.com/jawher/mow.cli"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	corePkg "github.com/opctl/opctl/cli/internal/core"
	"github.com/opctl/opctl/cli/internal/model"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	localNodeProvider "github.com/opctl/opctl/cli/internal/nodeprovider/local"
	op "github.com/opctl/opctl/sdks/go/opspec"
)

//counterfeiter:generate -o internal/fakes/cli.go . cli
type cli interface {
	Run(args []string) error
}

// newCorer allows swapping out corePkg.New for unit tests
type newCorer func(
	cliColorer clicolorer.CliColorer,
	nodeProvider nodeprovider.NodeProvider,
) corePkg.Core

func newCli(
	cliColorer clicolorer.CliColorer,
	newCorer newCorer,
) cli {

	cli := mow.App(
		"opctl",
		"Opctl is a free and open source distributed operation control system.",
	)
	cli.Version("v version", version)

	perUserAppDataPath, err := appdatapath.New().PerUser()
	if nil != err {
		panic(err)
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
	nodeCreateOpts := model.NodeCreateOpts{
		ContainerRuntime: *containerRuntime,
		DataDir:          *dataDir,
		ListenAddress:    *listenAddress,
	}

	core := newCorer(
		cliColorer,
		localNodeProvider.New(nodeCreateOpts),
	)

	noColor := cli.BoolOpt("nc no-color", false, "Disable output coloring")

	cli.Before = func() {
		if *noColor {
			cliColorer.Disable()
		}
	}

	cli.Command("auth", "Manage auth for OCI image registries", func(authCmd *mow.Cmd) {

		authCmd.Command(
			"add", "Add auth for an OCI image registry",
			func(addCmd *mow.Cmd) {
				addCmd.Spec = "RESOURCES [ -u=<username> ] [ -p=<password> ]"

				resources := addCmd.StringArg("RESOURCES", "", "Resources this auth applies to in the form of a host or host/path.")
				username := addCmd.StringOpt("u username", "", "Username")
				password := addCmd.StringOpt("p password", "", "Password")

				addCmd.Action = func() {
					core.Auth().Add(context.TODO(), *resources, *username, *password)
				}
			})

	},
	)

	cli.Command("events", "Stream events", func(eventsCmd *mow.Cmd) {
		eventsCmd.Action = func() {
			core.Events(context.TODO())
		}
	})

	cli.Command(
		"ls", "List operations (only valid ops will be listed)", func(lsCmd *mow.Cmd) {
			const dirRefArgName = "DIR_REF"
			lsCmd.Spec = fmt.Sprintf("[%v]", dirRefArgName)
			dirRef := lsCmd.StringArg(dirRefArgName, op.DotOpspecDirName, "Reference to dir ops will be listed from")
			lsCmd.Action = func() {
				core.Ls(context.TODO(), *dirRef)
			}
		})

	cli.Command("node", "Manage nodes", func(nodeCmd *mow.Cmd) {

		nodeCmd.Command("create", "Creates a node", func(createCmd *mow.Cmd) {
			createCmd.Action = func() {
				core.Node().Create(nodeCreateOpts)
			}
		})

		nodeCmd.Command("kill", "Kills a node", func(killCmd *mow.Cmd) {
			killCmd.Action = func() {
				core.Node().Kill()
			}
		})
	})

	cli.Command("op", "Manage ops", func(opCmd *mow.Cmd) {

		opCmd.Command(
			"create", "Create an op",
			func(createCmd *mow.Cmd) {
				path := createCmd.StringOpt("path", op.DotOpspecDirName, "Path the op will be created at")
				description := createCmd.StringOpt("d description", "", "Op description")
				name := createCmd.StringArg("NAME", "", "Op name")

				createCmd.Action = func() {
					core.Op().Create(*path, *description, *name)
				}
			})

		opCmd.Command(
			"install", "Install an op",
			func(installCmd *mow.Cmd) {
				path := installCmd.StringOpt("path", op.DotOpspecDirName, "Path the op will be installed at")
				opRef := installCmd.StringArg("OP_REF", "", "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)")
				username := installCmd.StringOpt("u username", "", "Username used to auth w/ the pkg source")
				password := installCmd.StringOpt("p password", "", "Password used to auth w/ the pkg source")

				installCmd.Action = func() {
					core.Op().Install(context.TODO(), *path, *opRef, *username, *password)
				}
			})

		opCmd.Command("kill", "Kill an op", func(killCmd *mow.Cmd) {
			opID := killCmd.StringArg("OP_ID", "", "Id of the op to kill")

			killCmd.Action = func() {
				core.Op().Kill(context.TODO(), *opID)
			}
		})

		opCmd.Command(
			"validate", "Validate an op",
			func(validateCmd *mow.Cmd) {
				opRef := validateCmd.StringArg("OP_REF", "", "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)")

				validateCmd.Action = func() {
					core.Op().Validate(context.TODO(), *opRef)
				}
			})

	})

	cli.Command("run", "Start and wait on an op", func(runCmd *mow.Cmd) {
		args := runCmd.StringsOpt("a", []string{}, "Explicitly pass args to op in format `-a NAME1=VALUE1 -a NAME2=VALUE2`")
		argFile := runCmd.StringOpt("arg-file", filepath.Join(op.DotOpspecDirName, "args.yml"), "Read in a file of args in yml format")
		opRef := runCmd.StringArg("OP_REF", "", "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)")

		runCmd.Action = func() {
			core.Run(context.TODO(), *opRef, &model.RunOpts{Args: *args, ArgFile: *argFile})
		}
	})

	cli.Command("self-update", "Update opctl", func(selfUpdateCmd *mow.Cmd) {
		channel := selfUpdateCmd.StringOpt("c channel", "stable", "Release channel to update from (either `stable`, `alpha`, or `beta`)")
		selfUpdateCmd.Action = func() {
			core.SelfUpdate(*channel)
		}
	})

	cli.Command("ui", "Open the opctl web UI and mount a reference.", func(uiCmd *mow.Cmd) {
		const mountRefArgName = "MOUNT_REF"
		uiCmd.Spec = fmt.Sprintf("[%v]", mountRefArgName)
		mountRefArg := uiCmd.StringArg(mountRefArgName, ".", "Reference to mount (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)")

		uiCmd.Action = func() {
			core.UI(*mountRefArg)
		}
	})

	return cli

}
