// Inspects a graph without passing secrets to plan.
package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Inputs struct {
	Token  opctl.String `opctl:"token,secret,default=PRIVATE-DEFAULT"`
	Region opctl.String `opctl:"region,default=west"`
}

func plan(in Inputs) opctl.Call {
	return opctl.Container("alpine:3").EnvFrom("TOKEN", in.Token).Env("REGION", in.Region.Value()).Cmd("printf", "THIS_CALL_MUST_NOT_RUN")
}
func main() { opctl.Main(plan) }
