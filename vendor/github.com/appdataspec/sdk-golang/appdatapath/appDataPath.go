package appdatapath

import "github.com/golang-interfaces/ios"

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ AppDataPath

type AppDataPath interface {
	// Global returns the per user app data path
	// returns non-nil error if required env vars missing
	Global() (string, error)
	// PerUser returns the per user app data path
	// returns non-nil error if required env vars missing
	PerUser() (string, error)
}

func New() AppDataPath {
	return appDataPath{
		os: ios.New(),
	}
}

type appDataPath struct {
	os ios.IOS
}
