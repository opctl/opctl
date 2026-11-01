package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Inputs struct {
	Default opctl.Object[map[string]string] `opctl:"default,default='{\"left\":\"from-handle\"}'"`
}

func main() {
	opctl.Main(func(in Inputs) opctl.Call {
		c := opctl.CallOp("github.com/opspec-pkgs/object.merge#1.0.0").Inputs(map[string]any{"default": in.Default, "override": map[string]string{"right": "$(not-a-ref)"}})
		c, result := opctl.OpOutput[opctl.Object[map[string]string]](c, "result")
		return opctl.Serial(c, opctl.Container("alpine:3").EnvFrom("RESULT", result).Env("EXPECTED", `{"left":"from-handle","right":"$(not-a-ref)"}`).Cmd("sh", "-ce", `test "$RESULT" = "$EXPECTED"`))
	})
}
