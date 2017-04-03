package vos

import (
	"os"
)

//go:generate counterfeiter -o ./fakeVOS.go --fake-name FakeVOS ./ VOS

// virtual operating system interface
type VOS interface {
	// Getenv retrieves the value of the environment variable named by the key.
	// It returns the value, which will be empty if the variable is not present.
	Getenv(key string) string

	// Setenv sets the value of the environment variable named by the key.
	// It returns an error, if any.
	Setenv(key, value string) error
}

func New() VOS {
	return _vos{}
}

type _vos struct{}

func (this _vos) Getenv(key string) string {
	return os.Getenv(key)
}

func (this _vos) Setenv(key, value string) error {
	return os.Setenv(key, value)
}
