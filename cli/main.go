package main

import (
	"fmt"
	"os"

	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/clioutput"
)

func main() {
	cliOutput := clioutput.New(clicolorer.New(), os.Stderr, os.Stdout)
	defer func() {
		if panicArg := recover(); panicArg != nil {
			cliOutput.Error(fmt.Sprint(panicArg))
			os.Exit(1)
		}
	}()

	newCli(
    cliOutput,
	).
		Run(os.Args)

}
