package core

//go:generate counterfeiter -o ./fakeExiter.go --fake-name fakeExiter ./ exiter

type ExitReq struct {
	Message string
	Code    int
}

// allows mocking/faking program exit
type exiter interface {
	Exit(req ExitReq)
}

func newExiter() exiter {
	return _exiter{}
}

type _exiter struct{}

func (this _exiter) Exit(req ExitReq) {
	panic(req)
}
