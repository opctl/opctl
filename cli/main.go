package main

import (
	"fmt"
	"github.com/opctl/opctl/cli/core"
	"github.com/opctl/opctl/util/clicolorer"
	"os"
)

func main() {

	_cliColorer := clicolorer.New()
	defer func() {
		if panicArg := recover(); panicArg != nil {
			fmt.Println(_cliColorer.Error("%v", panicArg))
			os.Exit(1)
		}
	}()

	newCli(core.New(_cliColorer), _cliColorer).Run(os.Args)

}
