package main

import (
	"fmt"
	"github.com/opspec-io/opctl/cli/core"
	"github.com/opspec-io/opctl/util/colorer"
	"os"
)

func main() {

	_colorer := colorer.New()
	defer func() {
		if panicArg := recover(); panicArg != nil {
			fmt.Println(_colorer.Error("%v", panicArg))
			os.Exit(1)
		}
	}()
	newCli(core.New(_colorer), _colorer).Run(os.Args)

}
