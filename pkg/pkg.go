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
	// Create creates a package
	Create(
		path,
		pkgName,
		pkgDescription string,
	) error

	// GetManifest gets the manifest of a package
	GetManifest(
		pkgHandle model.PkgHandle,
	) (
		*model.PkgManifest,
		error,
	)

	// Install installs a pkg; path will be created if it doesn't exist
	Install(
		ctx context.Context,
		path string,
		pkgHandle model.PkgHandle,
	) error

	// List recursively lists packages in dirPath
	List(
		dirPath string,
	) ([]*model.PkgManifest, error)

	providerFactory

	resolver

	// Validate validates a package
	Validate(
		pkgHandle model.PkgHandle,
	) []error
}

func New() Pkg {
	return _Pkg{
		ioUtil:          iioutil.New(),
		manifest:        manifest.New(),
		os:              ios.New(),
		providerFactory: newProviderFactory(),
		puller:          newPuller(),
		refParser:       newRefParser(),
		resolver:        newResolver(),
	}
}

type _Pkg struct {
	ioUtil   iioutil.IIOUtil
	manifest manifest.Manifest
	os       ios.IOS
	providerFactory
	puller
	refParser
	resolver
}
