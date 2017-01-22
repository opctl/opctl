package main

//go:generate counterfeiter -o ./fakeCli.go --fake-name FakeCli ./ cli

import (
	mow "github.com/jawher/mow.cli"
	"github.com/opspec-io/opctl/cli/core"
	"github.com/opspec-io/opctl/daemon"
	"github.com/opspec-io/opctl/util/colorer"
)

type cli interface {
	Run(args []string) error
}

func newCli(
	core core.Core,
	colorer colorer.Colorer,
) cli {

	cli := mow.App("opctl", "Cross platform runtime for https://opspec.io")
	cli.Version("v version", version)

	noColor := cli.BoolOpt("nc no-color", false, "Disable output coloring")
	cli.Before = func() {
		if *noColor {
			colorer.Disable()
		}
	}

	cli.Command("collection", "Collection related actions", func(collectionCmd *mow.Cmd) {

		collectionCmd.Command("create", "Create a collection", func(createCmd *mow.Cmd) {

			name := createCmd.StringArg("NAME", "", "Name of the collection")
			description := createCmd.StringOpt("d description", "", "Description of the collection")

			createCmd.Action = func() {
				core.CreateCollection(*description, *name)
			}

		})

		collectionCmd.Command("set", "Set collection attributes", func(setCmd *mow.Cmd) {
			setCmd.Command("description", "Set the description of a collection", func(descriptionCmd *mow.Cmd) {
				description := descriptionCmd.StringArg("DESCRIPTION", "", "Description of the collection")

				descriptionCmd.Action = func() {
					core.SetCollectionDescription(*description)
				}
			})
		})
	})

	cli.Command("daemon", "Run the opctl daemon", func(daemonCmd *mow.Cmd) {
		daemonCmd.Action = func() {
			daemon.New()
		}
	})

	cli.Command("events", "Stream events", func(eventsCmd *mow.Cmd) {
		eventsCmd.Action = func() {
			core.StreamEvents()
		}
	})

	cli.Command("kill", "Kill an op graph", func(killCmd *mow.Cmd) {
		opId := killCmd.StringArg("OP_GRAPH_ID", "", "Id of the op graph to kill")

		killCmd.Action = func() {
			core.KillOp(*opId)
		}
	})

	cli.Command("ls", "List ops in a collection", func(lsCmd *mow.Cmd) {
		collection := lsCmd.StringOpt("c collection", ".opspec", "Collection to list ops from")
		lsCmd.Action = func() {
			core.ListOpsInCollection(*collection)
		}
	})

	cli.Command("op", "Op related actions", func(opCmd *mow.Cmd) {

		opCmd.Command("create", "Create an op", func(createCmd *mow.Cmd) {
			collection := createCmd.StringOpt("c collection", ".opspec", "Collection to embed op in (use '' to not embed)")
			description := createCmd.StringOpt("d description", "", "Description of the op")
			name := createCmd.StringArg("NAME", "", "Name of the op")

			createCmd.Action = func() {
				core.CreateOp(*collection, *description, *name)
			}
		})

		opCmd.Command("set", "Set op attributes", func(setCmd *mow.Cmd) {
			setCmd.Command("description", "Set the description of an op", func(descriptionCmd *mow.Cmd) {
				collection := descriptionCmd.StringOpt("c collection", ".opspec", "Collection the op is embedded in (use '' if not embedded)")
				description := descriptionCmd.StringArg("DESCRIPTION", "", "description of the op")
				name := descriptionCmd.StringArg("NAME", "", "name of the op")

				descriptionCmd.Action = func() {
					core.SetOpDescription(*collection, *description, *name)
				}
			})
		})
	})

	cli.Command("run", "Run an op", func(runCmd *mow.Cmd) {
		args := runCmd.StringsOpt("a", []string{}, "Pass args to op in format: NAME[=VALUE] (gets from env if not provided)")
		collection := runCmd.StringOpt("c collection", ".opspec", "Collection the op is embedded in (use '' if not embedded)")
		name := runCmd.StringArg("OP_URL", "", "Url of the op (op name if in collection)")

		runCmd.Action = func() {
			core.RunOp(*args, *collection, *name)
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
