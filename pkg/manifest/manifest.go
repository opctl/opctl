// Package manifest implements use cases for managing opspec package manifests
package manifest

import (
	"github.com/golang-interfaces/iioutil"
	"github.com/opspec-io/sdk-golang/model"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Manifest

type Manifest interface {
	// Validate validates the pkg manifest at path
	Validate(path string) []error

	// Unmarshal unmarshals the pkg manifest at path
	Unmarshal(
		path string,
	) (*model.PkgManifest, error)
}

func New() Manifest {
	return _Manifest{
		validator: newValidator(),
		ioUtil:    iioutil.New(),
	}
}

type _Manifest struct {
	validator validator
	ioUtil    iioutil.Iioutil
}
