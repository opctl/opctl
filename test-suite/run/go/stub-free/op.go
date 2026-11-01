// This stub-free op mentions opctl:op in its documentation.
package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

func main() {
	opctl.Main(func(struct{}) opctl.Call { return opctl.Container("alpine:3").Cmd("true") })
}
