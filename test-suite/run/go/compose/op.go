package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

// greeting returns a call for the caller's graph and the directory it fills.
func greeting(who string) (opctl.Call, opctl.Dir) {
	out := opctl.NewDir()
	return opctl.Container("alpine:3").
		Cmd("sh", "-ce", `echo "hello $WHO" > "/out/$WHO.txt"`).
		Env("WHO", who).
		Dir("/out", out), out
}

func main() {
	opctl.Main(func(struct{}) opctl.Call {
		var greetings []opctl.Call
		combine := opctl.Container("alpine:3")
		for _, who := range []string{"a", "b"} {
			call, dir := greeting(who)
			greetings = append(greetings, call)
			combine = combine.Dir("/in/"+who, dir)
		}
		all := opctl.NewFile()
		combine = combine.Cmd("sh", "-ce", "cat /in/a/a.txt /in/b/b.txt > /all.txt").File("/all.txt", all)
		check := opctl.Container("alpine:3").
			File("/all.txt", all).
			File("/expected.txt", opctl.OpFile("./expected.txt")).
			Cmd("cmp", "/all.txt", "/expected.txt")
		return opctl.Serial(opctl.Parallel(greetings...), combine, check)
	})
}
