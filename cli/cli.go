package main

//go:generate counterfeiter -o ./fakeCli.go --fake-name FakeCli ./ cli

import (
	mow "github.com/jawher/mow.cli"
	"github.com/opctl/opctl/cli/core"
	"github.com/opctl/opctl/util/clicolorer"
)

type cli interface {
	Run(args []string) error
}

func newCli(
	core core.Core,
	cliColorer clicolorer.CliColorer,
) cli {

	cli := mow.App("opctl", "Open source cross platform opspec node")
	cli.Version("v version", version)

	noColor := cli.BoolOpt("nc no-color", false, "Disable output coloring")
	cli.Before = func() {
		if *noColor {
			cliColorer.Disable()
		}
	}

	cli.Command("events", "Stream events", func(eventsCmd *mow.Cmd) {
		eventsCmd.Action = func() {
			core.StreamEvents()
		}
	})

	cli.Command(
		"ls", "List packages (only v0.1.3 opspec packages will be listed)", func(lsCmd *mow.Cmd) {
			path := lsCmd.StringOpt("path", ".opspec", "Path to list packages from")
			lsCmd.Action = func() {
				core.ListPackages(*path)
			}
		})

	cli.Command("node", "Manage nodes", func(nodeCmd *mow.Cmd) {

		nodeCmd.Command("create", "Creates a node", func(createCmd *mow.Cmd) {
			createCmd.Action = func() {
				core.NodeCreate()
			}
		})

		nodeCmd.Command("kill", "Kills a node", func(killCmd *mow.Cmd) {
			killCmd.Action = func() {
				core.NodeKill()
			}
		})
	})

	cli.Command("op", "Manage ops", func(opCmd *mow.Cmd) {

		opCmd.Command("kill", "Kill an op", func(killCmd *mow.Cmd) {
			opId := killCmd.StringArg("OP_ID", "", "Id of the op to kill")

			killCmd.Action = func() {
				core.OpKill(*opId)
			}
		})

	})

	cli.Command("pkg", "Manage packages", func(pkgCmd *mow.Cmd) {

		pkgCmd.Command(
			"create", "Create a package",
			func(createCmd *mow.Cmd) {
				path := createCmd.StringOpt("path", ".opspec", "Path the package will be created at")
				description := createCmd.StringOpt("d description", "", "Package description")
				name := createCmd.StringArg("NAME", "", "Package name")

				createCmd.Action = func() {
					core.Create(*path, *description, *name)
				}
			})

		pkgCmd.Command("set", "Set attributes of a package", func(setCmd *mow.Cmd) {
			setCmd.Command("description", "Set the description of a package", func(descriptionCmd *mow.Cmd) {
				description := descriptionCmd.StringArg("DESCRIPTION", "", "Package description")
				pkgRef := descriptionCmd.StringArg("PACKAGE_REF", "", "Package reference")

				descriptionCmd.Action = func() {
					core.PkgSetDescription(*description, *pkgRef)
				}
			})
		})

		pkgCmd.Command(
			"validate", "Validates a package",
			func(validateCmd *mow.Cmd) {
				pkgRef := validateCmd.StringArg("PACKAGE_REF", "", "Package reference")

				validateCmd.Action = func() {
					core.PkgValidate(*pkgRef)
				}
			})
	})

	cli.Command("run", "Start and wait on an op", func(runCmd *mow.Cmd) {
		args := runCmd.StringsOpt("a", []string{}, "Pass args in format: NAME[=VALUE] (gets from env if not provided)")
		pkgRef := runCmd.StringArg("PACKAGE_REF", "", "Package reference")

		runCmd.Action = func() {
			core.RunOp(*args, *pkgRef)
		}
	})

	cli.Command("self-update", "Update opctl", func(selfUpdateCmd *mow.Cmd) {
		channel := selfUpdateCmd.StringOpt("c channel ", "stable", "Release channel to update from (either 'stable' or 'beta'; defaults to 'stable' )")
		selfUpdateCmd.Action = func() {
			core.SelfUpdate(*channel)
		}
	})

	return cli

}
