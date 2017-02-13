package core

//go:generate counterfeiter -o ./fakeExiter.go --fake-name fakeExiter ./ exiter

import (
	"github.com/opspec-io/opctl/util/vos"
)

type ExitReq struct {
	Message string
	Code    int
}

// allows mocking/faking program exit
type exiter interface {
	Exit(req ExitReq)
}

func newExiter(
	output output,
	vos vos.Vos,
) exiter {
	return _exiter{
		output: output,
		vos:    vos,
	}
}

type _exiter struct {
	output output
	vos    vos.Vos
}

func (this _exiter) Exit(req ExitReq) {
	if req.Code > 0 {
		this.output.Error(req.Message)
		this.vos.Exit(req.Code)
	} else {
		this.output.Success(req.Message)
		this.vos.Exit(0)
	}
}
