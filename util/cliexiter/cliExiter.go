package cliexiter

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ CliExiter

import (
	"github.com/opspec-io/opctl/util/clioutput"
	"github.com/opspec-io/opctl/util/vos"
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
	vos vos.Vos,
) CliExiter {
	return _cliExiter{
		cliOutput: cliOutput,
		vos:       vos,
	}
}

type _cliExiter struct {
	cliOutput clioutput.CliOutput
	vos       vos.Vos
}

func (this _cliExiter) Exit(req ExitReq) {
	if req.Code > 0 {
		this.cliOutput.Error(req.Message)
		this.vos.Exit(req.Code)
	} else {
		this.cliOutput.Success(req.Message)
		this.vos.Exit(0)
	}
}
