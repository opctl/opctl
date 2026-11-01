package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Outputs struct {
	Version  opctl.String                    `opctl:"version"`
	Count    opctl.Number                    `opctl:"count"`
	Ready    opctl.Boolean                   `opctl:"ready"`
	Targets  opctl.Array[string]             `opctl:"targets"`
	Metadata opctl.Object[map[string]string] `opctl:"metadata"`
}

func main() {
	opctl.MainWithOutputs(func(struct{}) (opctl.Call, Outputs) {
		producer := opctl.Container("alpine:3").Cmd("sh", "-ce", `printf '%s' '{"version":"v1.2.3","count":3,"ready":true,"targets":["linux"],"metadata":{"revision":"abc123"}}' > "$OPCTL_OUTPUTS"`)
		producer, version := opctl.ContainerOutput[opctl.String](producer, "version")
		producer, count := opctl.ContainerOutput[opctl.Number](producer, "count")
		producer, ready := opctl.ContainerOutput[opctl.Boolean](producer, "ready")
		producer, targets := opctl.ContainerOutput[opctl.Array[string]](producer, "targets")
		producer, metadata := opctl.ContainerOutput[opctl.Object[map[string]string]](producer, "metadata")
		check := opctl.Container("alpine:3").EnvFrom("VERSION", version).EnvFrom("COUNT", count).EnvFrom("READY", ready).EnvFrom("TARGETS", targets).EnvFrom("METADATA", metadata).
			Cmd("sh", "-ce", `test "$VERSION" = v1.2.3 && test "$COUNT" = 3 && test "$READY" = true && test "$TARGETS" = '["linux"]' && test "$METADATA" = '{"revision":"abc123"}'`)
		return opctl.Serial(producer, check), Outputs{version, count, ready, targets, metadata}
	})
}
