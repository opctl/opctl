// Package pkg implements use cases specific to opspec packages
package pkg

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Pkg

import (
	"context"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg/manifest"
)

type Pkg interface {
	ProviderFactory

	// Create creates a package
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

	// Install installs a pkg; path will be created if it doesn't exist
	Install(
		ctx context.Context,
		path string,
		pkgHandle Handle,
	) error

	// List recursively lists packages in dirPath
	List(
		dirPath string,
	) ([]*model.PkgManifest, error)

	Resolver

	// Validate validates a package
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
