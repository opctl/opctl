package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Metadata struct {
	Revision int `json:"revision"`
}
type Private struct {
	Token    opctl.String           `opctl:"token,secret"`
	Metadata opctl.Object[Metadata] `opctl:"metadata,secret,required=revision"`
}
type Public struct {
	Token    opctl.String           `opctl:"token"`
	Metadata opctl.Object[Metadata] `opctl:"metadata"`
}

func main() {
	opctl.MainWithOutputs(func(struct{}) (opctl.Call, Public) {
		producer, out := opctl.Outputs[Private](opctl.Container("alpine:3").Cmd("sh", "-ce", `token=$(printf '%032d' 42); printf '{"token":"%s","metadata":{"revision":"%s"}}' "$token" "$token" > "$OPCTL_OUTPUTS"`))
		return opctl.Serial(producer, opctl.Container("alpine:3").Cmd("sh", "-ce", "exit 99")), Public{Token: out.Token, Metadata: out.Metadata}
	})
}
