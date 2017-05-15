// Package pkg implements use cases for managing opspec packages
package pkg

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Pkg

import (
	"github.com/golang-interfaces/gopkg.in-src-d-go-git.v4"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg/manifest"
)

type Pkg interface {
	// Create creates an opspec package
	Create(
		path,
		pkgName,
		pkgDescription string,
	) error

	// Resolve attempts to resolve a package from lookPaths according to opspec package resolution rules.
	// if successful it's absolute path will be returned along w/ true
	Resolve(
		pkgRef *PkgRef,
		lookPaths ...string,
	) (string, bool)

	// ParseRef parses a pkgRef
	ParseRef(
		pkgRef string,
	) (*PkgRef, error)

	// Pull pulls 'pkgRef' to 'path'
	// returns ErrAuthenticationFailed on authentication failure
	Pull(
		path string,
		pkgRef *PkgRef,
		opts *PullOpts,
	) error

  // List recursively lists packages in dirPath
	List(
		dirPath string,
	) ([]*model.PkgManifest, error)

	// SetDescription sets the description of a package
	SetDescription(
		pkgPath,
		pkgDescription string,
	) error

	// Validate validates an opspec package
	Validate(
		pkgPath string,
	) []error
}

func New() Pkg {
	return _Pkg{
		git:      igit.New(),
		ioUtil:   iioutil.New(),
		os:       ios.New(),
		manifest: manifest.New(),
	}
}

type _Pkg struct {
	git      igit.IGit
	ioUtil   iioutil.Iioutil
	os       ios.IOS
	manifest manifest.Manifest
}
