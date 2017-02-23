package vruntime

import (
	"runtime"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Vruntime

// virtual operating system interface
type Vruntime interface {
	// GOOS is the running program's operating system target:
	// one of darwin, freebsd, linux, and so on.
	GOOS() string
}

func New() Vruntime {
	return _vruntime{}
}

type _vruntime struct{}

func (this _vruntime) GOOS() string {
	return runtime.GOOS
}
