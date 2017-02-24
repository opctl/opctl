package pscanary

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ PsCanary

type PsCanary interface {
	IsAlive(processId int) bool
}

func New() PsCanary {
	return psCanary{}
}

type psCanary struct{}
