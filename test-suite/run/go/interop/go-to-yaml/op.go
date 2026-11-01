package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Inputs struct {
	Src   opctl.Dir    `opctl:"src,default=$(./src)"`
	Label opctl.String `opctl:"label,default=caller-label"`
	Mode  opctl.String `opctl:"mode,default=normal"`
}

func main() {
	opctl.Main(func(in Inputs) opctl.Call {
		mode := in.Mode.Value()
		label := "from-go"
		if mode == "empty" {
			label = ""
		}
		args := map[string]any{"src": in.Src, "label": label, "payload": map[string]any{"$key": "$(not-a-ref)"}}
		if mode == "unknown-input" {
			args["unknown"] = true
		}
		c := opctl.CallOp("./yaml").Inputs(args)
		c, artifact := opctl.OpOutput[opctl.Dir](c, "artifact")
		if mode == "incompatible" {
			c, bad := opctl.OpOutput[opctl.Dir](c, "label")
			return opctl.Serial(c, opctl.Container("alpine:3").Dir("/bad", bad))
		}
		name := "label"
		if mode == "unknown-output" {
			name = "lable"
		}
		c, echoed := opctl.OpOutput[opctl.String](c, name)
		check := opctl.Container("alpine:3").Dir("/artifact", artifact).Dir("/src", in.Src).EnvFrom("LABEL", echoed).Env("EXPECTED", label).Cmd("sh", "-ce", `cmp /artifact/marker /src/marker; test "$LABEL" = "$EXPECTED"; test "$(cat /artifact/label)" = "$EXPECTED"`)
		if mode == "parallel" {
			return opctl.Parallel(c, check)
		}
		return opctl.Serial(c, check)
	})
}
