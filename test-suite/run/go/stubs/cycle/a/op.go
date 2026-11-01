package main

import "opctl.test/go-stubs/opctlgen/cycleb" //opctl:op ../b
import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

func main() { opctl.Main(func(struct{}) opctl.Call { c, _ := cycleb.Call(cycleb.Inputs{}); return c }) }
