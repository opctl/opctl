package main

//go:generate counterfeiter -o ./fakeCli.go --fake-name FakeCli ./ cli

import (
	"context"
	"fmt"
	"path/filepath"

	mow "github.com/jawher/mow.cli"
	nodeCreateCmd "github.com/opctl/opctl/cli/cmds/node/create"
	corePkg "github.com/opctl/opctl/cli/core"
	"github.com/opctl/opctl/cli/model"
	"github.com/opctl/opctl/cli/util/clicolorer"
	"github.com/opctl/sdk-golang/opspec"
)

type cli interface {
	Run(args []string) error
}

func newCli(
	core corePkg.Core,
	cliColorer clicolorer.CliColorer,
	nodeCreateCmdInvoker nodeCreateCmd.Invoker,
) cli {

	cli := mow.App("opctl", "Opctl is a free and open source distributed operation control system.")
	cli.Version("v version", version)

	noColor := cli.BoolOpt("nc no-color", false, "Disable output coloring")
	cli.Before = func() {
		if *noColor {
			cliColorer.Disable()
		}
	}

	cli.Command("events", "Stream events", func(eventsCmd *mow.Cmd) {
		eventsCmd.Action = func() {
			core.Events(context.TODO())
		}
	})

	cli.Command(
		"ls", "List operations (only opspec 0.1.6 compatible operations will be listed)", func(lsCmd *mow.Cmd) {
			const DirRefArgName = "DIR_REF"
			lsCmd.Spec = fmt.Sprintf("[%v]", DirRefArgName)
			dirRef := lsCmd.StringArg(DirRefArgName, op.DotOpspecDirName, "Reference to dir ops will be listed from")
			lsCmd.Action = func() {
				core.Ls(context.TODO(), *dirRef)
			}
		})

	cli.Command("node", "Manage nodes", func(nodeCmd *mow.Cmd) {

		nodeCmd.Command("create", "Creates a node", func(createCmd *mow.Cmd) {
			dataDir := createCmd.StringOpt("data-dir", "", "Path of dir used to store node data (since v0.1.25)")
			if "" == *dataDir {
				dataDir = nil
			}

			createCmd.Action = func() {
				nodeCreateCmdInvoker.Invoke(
					model.NodeCreateOpts{
						DataDir: dataDir,
					},
				)
			}
		})

		nodeCmd.Command("kill", "Kills a node", func(killCmd *mow.Cmd) {
			killCmd.Action = func() {
				core.NodeKill()
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
					core.OpCreate(*path, *description, *name)
				}
			})

		opCmd.Command(
			"install", "Installs an op",
			func(installCmd *mow.Cmd) {
				defaultPath := fmt.Sprintf("%v/OP_REF", op.DotOpspecDirName)
				path := installCmd.StringOpt("path", defaultPath, "Path the op will be installed at")
				opRef := installCmd.StringArg("OP_REF", "", "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag` (since v0.1.19), or `host/path/repo#tag/path` (since v0.1.24))")
				username := installCmd.StringOpt("u username", "", "Username used to auth w/ the pkg source")
				password := installCmd.StringOpt("p password", "", "Password used to auth w/ the pkg source")

				installCmd.Action = func() {
					if *path == defaultPath {
						*path = filepath.Join(op.DotOpspecDirName, *opRef)
					}
					core.OpInstall(context.TODO(), *path, *opRef, *username, *password)
				}
			})

		opCmd.Command("kill", "Kill an op", func(killCmd *mow.Cmd) {
			opId := killCmd.StringArg("OP_ID", "", "Id of the op to kill")

			killCmd.Action = func() {
				core.OpKill(context.TODO(), *opId)
			}
		})

		opCmd.Command(
			"validate", "Validates an op",
			func(validateCmd *mow.Cmd) {
				opRef := validateCmd.StringArg("OP_REF", "", "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag` (since v0.1.19), or `host/path/repo#tag/path` (since v0.1.24))")

				validateCmd.Action = func() {
					core.OpValidate(context.TODO(), *opRef)
				}
			})

	})

	cli.Command("run", "Start and wait on an op", func(runCmd *mow.Cmd) {
		args := runCmd.StringsOpt("a", []string{}, "Explicitly pass args to op in format `-a NAME1=VALUE1 -a NAME2=VALUE2`")
		argFile := runCmd.StringOpt("arg-file", filepath.Join(op.DotOpspecDirName, "args.yml"), "Read in a file of args in yml format")
		opRef := runCmd.StringArg("OP_REF", "", "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag` (since v0.1.19), or `host/path/repo#tag/path` (since v0.1.24))")

		runCmd.Action = func() {
			core.Run(context.TODO(), *opRef, &model.RunOpts{Args: *args, ArgFile: *argFile})
		}
	})

	cli.Command("self-update", "Update opctl", func(selfUpdateCmd *mow.Cmd) {
		channel := selfUpdateCmd.StringOpt("c channel ", "stable", "Release channel to update from (either `stable`, `alpha`, or `beta`)")
		selfUpdateCmd.Action = func() {
			core.SelfUpdate(*channel)
		}
	})

	return cli

}
