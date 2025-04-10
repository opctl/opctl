package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/opctl/opctl/cli/cmd"
)

func main() {
	defer func() {
		if panic := recover(); panic != nil {
			fmt.Fprintf(
				os.Stderr,
				fmt.Sprintf("recovered from panic: %s\n%s", panic, string(debug.Stack())),
			)
			os.Exit(1)
		}
	}()

	cmd.Execute()

}
