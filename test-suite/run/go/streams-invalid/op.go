package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

func main() {
	opctl.Main(func(struct{}) opctl.Call {
		producer, version := opctl.Container("alpine:3").Cmd("echo", "v1.2.3").Stdout()
		return opctl.Parallel(producer, opctl.Container("alpine:3").EnvFrom("VERSION", version).Cmd("true"))
	})
}
