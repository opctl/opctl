package main

import (
	opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"
	"opctl.test/go-stubs/opctlgen/gogreet" //opctl:op ../gogreet
	"opctl.test/go-stubs/opctlgen/merge"   //opctl:op github.com/opspec-pkgs/object.merge#1.0.0
	"opctl.test/go-stubs/opctlgen/yaml"    //opctl:op ../yaml
)

type Outputs struct {
	Greeting opctl.String `opctl:"greeting"`
}

func main() {
	opctl.MainWithOutputs(func(struct{}) (opctl.Call, Outputs) {
		y, yout := yaml.Call(yaml.Inputs{Src: opctl.OpDir("./assets"), Label: opctl.StringOf("from-go"), Tags: opctl.ArrayOf([]string{"linux"}), Config: opctl.ObjectOf(yaml.InputsConfig{Style: "$(not-a-ref)"})})
		g, gout := gogreet.Call(gogreet.Inputs{Name: opctl.StringOf("from-go")})
		m, mout := merge.Call(merge.Inputs{Default: opctl.ObjectOf(map[string]any{"left": "local"}), Override: opctl.ObjectOf(map[string]any{"right": "remote"})})
		check := opctl.Container("alpine:3").Dir("/artifact", yout.Artifact).EnvFrom("TEXT", yout.Text).EnvFrom("GREETING", gout.Greeting).EnvFrom("MERGED", mout.Result).
			Cmd("sh", "-ce", `test "$(cat /artifact/marker)" = 'typed op calls'; test "$TEXT" = default-from-go; test "$GREETING" = 'hello from-go'; test "$MERGED" = '{"left":"local","right":"remote"}'`)
		return opctl.Serial(y, g, m, check), Outputs{Greeting: gout.Greeting}
	})
}
