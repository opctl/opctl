package main

import (
	"example.com/app/internal/greeting"
	opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"
)

type Inputs struct {
	Suffix opctl.String `opctl:"suffix,default=first"`
}

func main() {
	opctl.Main(func(in Inputs) opctl.Call {
		return opctl.Container("alpine:3").
			Env("GREETING", greeting.Text()).EnvFrom("SUFFIX", in.Suffix).
			Cmd("sh", "-ce", `test "$GREETING" = 'hello from the app module'; case "$SUFFIX" in first|second) echo "$GREETING $SUFFIX";; *) exit 1;; esac`)
	})
}
