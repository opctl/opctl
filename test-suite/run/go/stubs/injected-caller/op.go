package main

import (
	"opctl.test/go-stubs/opctlgen/injected" //opctl:op ../injected-callee

	opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"
)

func main() {
	opctl.Main(func(struct{}) opctl.Call {
		call, _ := injected.Call(injected.Inputs{})
		return call
	})
}
