package sockets

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Sockets

import (
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
)

type Sockets interface {
	Interpret(
		scope map[string]*model.Value,
		scgContainerCallSockets map[string]string,
		scratchDirPath string,
	) (map[string]string, error)
}

func New() Sockets {
	return _Sockets{
		os: ios.New(),
	}
}

type _Sockets struct {
	os ios.IOS
}
