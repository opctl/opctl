package pscanary

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ PsCanary

// allows mocking/faking program exit
type PsCanary interface {
	IsAlive(processId int) bool
}

func New() PsCanary {
	return psCanary{}
}

type psCanary struct{}
