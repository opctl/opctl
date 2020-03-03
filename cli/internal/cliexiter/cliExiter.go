package cliexiter

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/cli/internal/clioutput"
)

type ExitReq struct {
	Message string
	Code    int
}

//CliExiter allows mocking/faking program exit
//counterfeiter:generate -o fakes/cliExiter.go . CliExiter
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
