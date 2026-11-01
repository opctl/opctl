package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Inputs struct {
	Marker opctl.String `opctl:"marker,default=schema-one"`
}

func main() {
	opctl.Main(func(in Inputs) opctl.Call {
		file := opctl.NewFile()
		return opctl.Serial(
			opctl.Container("alpine:3").File("/marker", file).EnvFrom("MARKER", in.Marker).Cmd("sh", "-ce", `printf '%s' "$MARKER" > /marker`),
			opctl.Container("alpine:3").File("/marker", file).EnvFrom("MARKER", in.Marker).Cmd("sh", "-ce", `test "$(cat /marker)" = "$MARKER"`),
		)
	})
}
