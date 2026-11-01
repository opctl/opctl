package main

import "opctl.test/go-stubs/opctlgen/cyclea" //opctl:op ../a
import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

func main() { opctl.Main(func(struct{}) opctl.Call { c, _ := cyclea.Call(cyclea.Inputs{}); return c }) }
