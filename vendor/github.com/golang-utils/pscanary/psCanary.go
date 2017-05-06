package pscanary

import "github.com/golang-interfaces/vos"

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ PsCanary

type PsCanary interface {
	IsAlive(processId int) bool
}

func New() PsCanary {
	return _PsCanary{
		os: vos.New(),
	}
}

type _PsCanary struct {
	os vos.VOS
}
