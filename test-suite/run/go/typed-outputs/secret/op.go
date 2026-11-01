package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Metadata struct {
	Revision string `json:"revision"`
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
		relay := opctl.CallOp("../../../container/object/typed-outputs/secret/relay").Inputs(map[string]any{"token": out.Token, "metadata": out.Metadata})
		relay, renamed := opctl.OpOutput[opctl.String](relay, "token")
		relay, details := opctl.OpOutput[opctl.Object[Metadata]](relay, "metadata")
		check := func() opctl.ContainerCall {
			return opctl.Container("alpine:3").EnvFrom("TOKEN", renamed).EnvFrom("METADATA", details).Cmd("sh", "-ce", `expected=$(printf '%032d' 42); test "$TOKEN" = "$expected"; test "$METADATA" = "{\"revision\":\"$expected\"}"`)
		}
		return opctl.Serial(opctl.Parallel(producer, opctl.Container("alpine:3").Cmd("true")), relay, opctl.Parallel(check(), check()), check()), Public{Token: out.Token, Metadata: out.Metadata}
	})
}
