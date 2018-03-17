// Package pkg implements use cases specific to ops
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
	// Create creates an operation
	Create(
		path,
		pkgName,
		pkgDescription string,
	) error

	// GetManifest gets the manifest of an operation
	GetManifest(
		opDirHandle model.DataHandle,
	) (
		*model.PkgManifest,
		error,
	)

	// Install installs an op; path will be created if it doesn't exist
	Install(
		ctx context.Context,
		path string,
		opDirHandle model.DataHandle,
	) error

	// List recursively lists ops in dirPath
	ListOps(
		dirPath string,
	) ([]*model.PkgManifest, error)

	// Validate validates an op
	Validate(
		opDirHandle model.DataHandle,
	) []error
}

func New() Pkg {
	return _Pkg{
		ioUtil:   iioutil.New(),
		manifest: manifest.New(),
		os:       ios.New(),
	}
}

type _Pkg struct {
	ioUtil   iioutil.IIOUtil
	manifest manifest.Manifest
	os       ios.IOS
}
