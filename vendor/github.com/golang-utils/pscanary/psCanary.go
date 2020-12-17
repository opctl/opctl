package pscanary

import "github.com/golang-interfaces/ios"

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ PsCanary

type PsCanary interface {
	IsAlive(processId int) bool
}

func New() PsCanary {
	return _PsCanary{
		os: ios.New(),
	}
}

type _PsCanary struct {
	os ios.IOS
}
