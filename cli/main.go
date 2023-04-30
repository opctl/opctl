package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/clioutput"
)

func main() {
	cliOutput := clioutput.New(clicolorer.New(), os.Stderr, os.Stdout)
	defer func() {
		if panic := recover(); panic != nil {
			cliOutput.Error(
				fmt.Sprintf("recovered from panic: %s\n%s", panic, string(debug.Stack())),
			)
			os.Exit(1)
		}
	}()

	newCli(
		cliOutput,
	).
		Run(os.Args)

}
