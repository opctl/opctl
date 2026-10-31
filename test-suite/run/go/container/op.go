package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

func main() {
	opctl.Main(func(struct{}) opctl.Call {
		return opctl.Container("alpine:3").Cmd("echo", "hello from a Go op").Env("A", "b").WorkDir("/")
	})
}
