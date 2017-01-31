package path

import (
	vosPkg "github.com/appdataspec/sdk-golang/util/vos"
)

//go:generate counterfeiter -o ./fakePath.go --fake-name FakePath ./ Path

type Path interface {
	// returns the per user app data path; panics if required env vars missing
	Global() string
	// returns the per user app data path; panics if required env vars missing
	PerUser() string
}

func New() Path {
	return NewWithVos(vosPkg.New())
}

// allows passing in virtual operating system to decouple from running OS (useful for testing)
func NewWithVos(
	vos vosPkg.Vos,
) Path {
	return path{
		vos: vos,
	}
}

type path struct {
	vos vosPkg.Vos
}
