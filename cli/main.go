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
			switch arg := panicArg.(type) {
			case core.ExitReq:
				if arg.Code > 0 {
					fmt.Println(_colorer.Error(arg.Message))
				} else {
					fmt.Println(_colorer.Success(arg.Message))
				}
				os.Exit(arg.Code)
			default:
				fmt.Println(_colorer.Error("%v", arg))
				os.Exit(1)
			}
		}
	}()
	newCli(core.New(_colorer), _colorer).Run(os.Args)

}
