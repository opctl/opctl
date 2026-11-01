package main

import "opctl.test/go-stubs/opctlgen/nope"
import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

func main() { opctl.Main(func(struct{}) opctl.Call { c, _ := nope.Call(nope.Inputs{}); return c }) }
