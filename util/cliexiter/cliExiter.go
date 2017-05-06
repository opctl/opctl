package cliexiter

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ CliExiter

import (
	"github.com/golang-interfaces/vos"
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
	vos vos.VOS,
) CliExiter {
	return cliExiter{
		cliOutput: cliOutput,
		vos:       vos,
	}
}

type cliExiter struct {
	cliOutput clioutput.CliOutput
	vos       vos.VOS
}

func (this cliExiter) Exit(req ExitReq) {
	if req.Code > 0 {
		this.cliOutput.Error(req.Message)
		this.vos.Exit(req.Code)
	} else {
		this.cliOutput.Success(req.Message)
		this.vos.Exit(0)
	}
}
