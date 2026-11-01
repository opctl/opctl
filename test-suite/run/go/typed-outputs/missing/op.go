package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Results struct {
	Version  opctl.String        `opctl:"version"`
	Count    opctl.Number        `opctl:"count,minimum=1,format=integer"`
	Ready    opctl.Boolean       `opctl:"ready"`
	Targets  opctl.Array[string] `opctl:"targets,minItems=1"`
	Metadata opctl.Object[struct {
		Revision string `json:"revision"`
	}] `opctl:"metadata,required=revision"`
}

func main() {
	opctl.MainWithOutputs(func(struct{}) (opctl.Call, Results) {
		producer, out := opctl.Outputs[Results](opctl.Container("alpine:3").Cmd("sh", "-ce", `printf 'v1.2.3\n'; printf 'diagnostic\n' >&2; printf '%s' '{"version":"v1.2.3","ready":true,"targets":["linux"],"metadata":{"revision":"abc123"}}' > "$OPCTL_OUTPUTS"`))
		producer, stream := producer.Stdout()
		producer, diagnostic := producer.Stderr()
		check := opctl.Container("alpine:3").EnvFrom("VERSION", out.Version).EnvFrom("COUNT", out.Count).EnvFrom("READY", out.Ready).EnvFrom("TARGETS", out.Targets).EnvFrom("METADATA", out.Metadata).EnvFrom("STREAM", stream).EnvFrom("DIAGNOSTIC", diagnostic).
			Cmd("sh", "-ce", `test "$VERSION" = v1.2.3 && test "$COUNT" = 3 && test "$READY" = true && test "$TARGETS" = '["linux"]' && test "$METADATA" = '{"revision":"abc123"}' && test "$STREAM" = v1.2.3 && test "$DIAGNOSTIC" = diagnostic`)
		return opctl.Serial(producer, check), out
	})
}
