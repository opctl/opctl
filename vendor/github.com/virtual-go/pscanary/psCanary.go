package pscanary

import "github.com/virtual-go/vos"

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ PsCanary

type PsCanary interface {
	IsAlive(processId int) bool
}

func New(
	os vos.VOS,
) PsCanary {
	return _PsCanary{
		os: os,
	}
}

type _PsCanary struct {
	os vos.VOS
}
