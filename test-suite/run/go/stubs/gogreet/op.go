package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Inputs struct {
	Name opctl.String `opctl:"name"`
}
type Outputs struct {
	Greeting opctl.String `opctl:"greeting"`
}

func main() {
	opctl.MainWithOutputs(func(in Inputs) (opctl.Call, Outputs) {
		c, out := opctl.Outputs[Outputs](opctl.Container("alpine:3").EnvFrom("NAME", in.Name).Cmd("sh", "-ce", `printf '{"greeting":"hello %s"}' "$NAME" > "$OPCTL_OUTPUTS"`))
		return c, out
	})
}
