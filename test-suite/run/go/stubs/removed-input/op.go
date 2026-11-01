package main

import "opctl.test/go-stubs/opctlgen/yaml" //opctl:op ../yaml
import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

func main() {
	opctl.Main(func(struct{}) opctl.Call { c, _ := yaml.Call(yaml.Inputs{Gone: opctl.StringOf("gone")}); return c })
}
