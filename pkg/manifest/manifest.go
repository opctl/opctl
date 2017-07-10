// Package manifest implements use cases for managing opspec package manifests
package manifest

import (
	"github.com/golang-interfaces/iioutil"
	"github.com/opspec-io/sdk-golang/model"
	"io"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Manifest

type Manifest interface {
	Validator

	// Unmarshal unmarshals the pkg manifest at path
	Unmarshal(
		manifestReader io.Reader,
	) (*model.PkgManifest, error)
}

func New() Manifest {
	return _Manifest{
		Validator: newValidator(),
		ioUtil:    iioutil.New(),
	}
}

type _Manifest struct {
	Validator
	ioUtil iioutil.IIOUtil
}
