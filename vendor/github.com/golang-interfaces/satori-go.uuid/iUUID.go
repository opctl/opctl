package iuuid

//go:generate counterfeiter -o fake.go --fake-name Fake ./ IUUID

import "github.com/satori/go.uuid"

type IUUID interface {
	// NewV4 returns random generated UUID.
	NewV4() uuid.UUID
}

func New() IUUID {
	return _IUUID{}
}

type _IUUID struct{}

func (_iuuid _IUUID) NewV4() uuid.UUID {
	return uuid.NewV4()
}
