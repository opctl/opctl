package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

func main() {
	opctl.Main(func(struct{}) opctl.Call {
		out := opctl.NewDir()
		producer := opctl.Container("alpine:3").Cmd("touch", "/out/marker").Dir("/out", out)
		consumer := opctl.Container("alpine:3").Cmd("test", "-f", "/in/marker").Dir("/in", out)
		return opctl.Parallel(producer, consumer)
	})
}
