package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Outputs struct {
	Token opctl.String `opctl:"token,secret"`
}

func main() {
	opctl.MainWithOutputs(func(struct{}) (opctl.Call, Outputs) {
		producer, token := opctl.Container("alpine:3").Cmd("sh", "-ce", `head -c 24 /dev/urandom | base64`).Stdout()
		checker := func() opctl.Call {
			return opctl.Container("alpine:3").EnvFrom("TOKEN", token).Cmd("sh", "-ce", `test "${#TOKEN}" = 32`)
		}
		return opctl.Serial(producer, opctl.Parallel(checker(), checker())), Outputs{token}
	})
}
