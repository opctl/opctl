// Package pkg implements use cases for managing opspec packages
package pkg

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Pkg

import (
	"github.com/golang-interfaces/gopkg.in-src-d-go-git.v4"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
)

type Pkg interface {
	// Create creates an opspec package
	Create(
		path,
		pkgName,
		pkgDescription string,
	) error

	// Resolve resolves a local package according to opspec package resolution rules and returns it's absolute path.
	Resolve(
		basePath,
		pkgRef string,
	) (string, bool)

	// Pull pulls a package from a remote source
	// returns ErrAuthenticationFailed on authentication failure
	Pull(
		pkgRef string,
		opts *PullOpts,
	) error

	// Get gets a local package
	Get(
		pkgPath string,
	) (*model.PkgManifest, error)

	// List lists packages according to opspec package resolution rules
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
	ioUtil := iioutil.New()
	manifestValidator := newManifestValidator()
	manifestUnmarshaller := newManifestUnmarshaller(ioUtil, manifestValidator)

	return pkg{
		git:                  igit.New(),
		ioUtil:               ioUtil,
		os:                   ios.New(),
		manifestUnmarshaller: manifestUnmarshaller,
		manifestValidator:    manifestValidator,
	}
}

type pkg struct {
	git                  igit.IGit
	ioUtil               iioutil.Iioutil
	os                   ios.IOS
	manifestValidator    manifestValidator
	manifestUnmarshaller manifestUnmarshaller
}
