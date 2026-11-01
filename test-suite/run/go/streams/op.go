package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

func main() {
	opctl.Main(func(struct{}) opctl.Call {
		producer, version := opctl.Container("alpine:3").Cmd("sh", "-ce", `printf 'v1.2.3\n'; printf 'diagnostic\n' >&2`).Stdout()
		producer, diagnostic := producer.Stderr()
		return opctl.Serial(producer, opctl.Container("alpine:3").EnvFrom("VERSION", version).EnvFrom("DIAGNOSTIC", diagnostic).Cmd("sh", "-ce", `test "$VERSION" = v1.2.3 && test "$DIAGNOSTIC" = diagnostic`))
	})
}
