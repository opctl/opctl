package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Inputs struct {
	Text  opctl.String `opctl:"text,default=from-cli"`
	Count opctl.Number `opctl:"count,default=7"`
	Token opctl.String `opctl:"token,secret,default=root-private-token"`
}
type Outputs struct {
	Echoed   opctl.String `opctl:"echoed"`
	Count    opctl.Number `opctl:"count"`
	Fallback opctl.String `opctl:"fallback,default=from-default"`
	Private  opctl.String `opctl:"private,secret"`
}

func main() {
	opctl.MainWithOutputs(func(in Inputs) (opctl.Call, Outputs) {
		return opctl.Container("alpine:3").Cmd("true"), Outputs{Echoed: in.Text, Count: in.Count, Private: in.Token}
	})
}
