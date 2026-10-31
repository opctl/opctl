package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Inputs struct {
	X opctl.Number `opctl:"x,minimun=1"`
}

func main() {
	opctl.Main(func(Inputs) opctl.Call { return opctl.Container("alpine:3") })
}
