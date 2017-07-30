// Package pkg implements use cases specific to opspec packages
package pkg

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Pkg

import (
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg/manifest"
)

type Pkg interface {
	ProviderFactory

	// Create creates an opspec package
	Create(
		path,
		pkgName,
		pkgDescription string,
	) error

	// GetManifest gets the manifest of a package
	GetManifest(
		pkgHandle Handle,
	) (
		*model.PkgManifest,
		error,
	)

	// List recursively lists packages in dirPath
	List(
		dirPath string,
	) ([]*model.PkgManifest, error)

	Resolver

	// Validate validates an opspec package
	Validate(
		pkgHandle Handle,
	) []error
}

func New() Pkg {
	return _Pkg{
		ioUtil:          iioutil.New(),
		manifest:        manifest.New(),
		os:              ios.New(),
		puller:          newPuller(),
		refParser:       newRefParser(),
		Resolver:        newResolver(),
		ProviderFactory: newProviderFactory(),
	}
}

type _Pkg struct {
	ioUtil   iioutil.IIOUtil
	manifest manifest.Manifest
	os       ios.IOS
	puller
	refParser
	Resolver
	ProviderFactory
}
