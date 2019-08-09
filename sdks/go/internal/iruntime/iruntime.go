package iruntime

import (
	"runtime"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fake.go --fake-name Fake ./ IRuntime

// virtual operating system interface
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
