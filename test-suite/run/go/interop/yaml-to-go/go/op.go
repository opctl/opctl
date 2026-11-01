package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Inputs struct {
	Src  opctl.Dir    `opctl:"src"`
	Text opctl.String `opctl:"text,minLength=3"`
	Kind opctl.String `opctl:"kind,default=container"`
}
type Outputs struct {
	Artifact opctl.Dir    `opctl:"artifact,description='Built files'"`
	Echoed   opctl.String `opctl:"echoed,minLength=3"`
	Source   opctl.Dir    `opctl:"source"`
}

func main() {
	opctl.MainWithOutputs(func(in Inputs) (opctl.Call, Outputs) {
		if in.Kind.Value() == "op" {
			c := opctl.CallOp("./identity").Inputs(map[string]any{"src": in.Src, "text": in.Text})
			c, artifact := opctl.OpOutput[opctl.Dir](c, "src")
			return c, Outputs{Artifact: artifact, Echoed: in.Text, Source: in.Src}
		}
		artifact := opctl.NewDir()
		c := opctl.Container("alpine:3").Dir("/src", in.Src).Dir("/out", artifact).EnvFrom("TEXT", in.Text).Cmd("sh", "-ce", `cp /src/marker /out/marker; printf '%s' "$TEXT" > /out/text`)
		return c, Outputs{Artifact: artifact, Echoed: in.Text, Source: in.Src}
	})
}
