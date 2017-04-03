package appdatapath

import (
	vos "github.com/appdataspec/sdk-golang/util/vos"
)

//go:generate counterfeiter -o ./fakeAppDataPath.go --fake-name FakePath ./ Path

type AppDataPath interface {
	// returns the per user app data path; panics if required env vars missing
	Global() string
	// returns the per user app data path; panics if required env vars missing
	PerUser() string
}

func New() AppDataPath {
	return NewWithVos(vos.New())
}

// allows passing in virtual operating system to decouple from running OS (useful for testing)
func NewWithVos(
	os vos.Vos,
) AppDataPath {
	return appDataPath{
		os: os,
	}
}

type appDataPath struct {
	os vos.Vos
}
