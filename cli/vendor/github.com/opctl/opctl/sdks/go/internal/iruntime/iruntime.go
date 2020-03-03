package iruntime

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"runtime"
)

// virtual operating system interface
//counterfeiter:generate -o fakes/iruntime.go . IRuntime
type IRuntime interface {
	// GOOS is the running program's operating system target:
	// one of darwin, freebsd, linux, and so on.
	GOOS() string
}

func New() IRuntime {
	return _IRuntime{}
}

type _IRuntime struct{}

func (this _IRuntime) GOOS() string {
	return runtime.GOOS
}
