package cliexiter

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ CliExiter

import (
	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/util/clioutput"
)

type ExitReq struct {
	Message string
	Code    int
}

// allows mocking/faking program exit
type CliExiter interface {
	Exit(req ExitReq)
}

func New(
	cliOutput clioutput.CliOutput,
	ios ios.IOS,
) CliExiter {
	return cliExiter{
		cliOutput: cliOutput,
		ios:       ios,
	}
}

type cliExiter struct {
	cliOutput clioutput.CliOutput
	ios       ios.IOS
}

func (this cliExiter) Exit(req ExitReq) {
	if req.Code > 0 {
		this.cliOutput.Error(req.Message)
		this.ios.Exit(req.Code)
	} else {
		this.cliOutput.Success(req.Message)
		this.ios.Exit(0)
	}
}
