package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Outputs struct {
	Count opctl.Number `opctl:"count,minimum=1,format=integer"`
}

func main() {
	opctl.MainWithOutputs(func(struct{}) (opctl.Call, Outputs) {
		producer, out := opctl.Outputs[Outputs](opctl.Container("alpine:3").Cmd("sh", "-ce", `printf '{"count":3}' > "$OPCTL_OUTPUTS"`))
		consumer := opctl.Container("alpine:3").EnvFrom("COUNT", out.Count).Cmd("sh", "-ce", `test "$COUNT" = 3`)
		return opctl.Serial(producer, consumer), out
	})
}
