// Package pkg implements use cases for managing opspec packages
package pkg

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Pkg

import (
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

	// ListContents lists contents of a package
	ListContents(
		pkgRef string,
	) (
		[]*model.PkgContent,
		error,
	)

	// GetContent gets content from a package
	GetContent(
		pkgRef string,
		contentPath string,
	) (
		model.ReadSeekCloser,
		error,
	)

	// Validate validates an opspec package
	Validate(
		pkgPath string,
	) []error
}

func New(
	cachePath string,
) Pkg {
	return _Pkg{
		ioUtil:    iioutil.New(),
		os:        ios.New(),
		puller:    newPuller(),
		refParser: newRefParser(),
		resolver:  newResolver(),
		manifest:  manifest.New(),
		opener:    newOpener(cachePath),
		cachePath: cachePath,
	}
}

type _Pkg struct {
	ioUtil iioutil.Iioutil
	os     ios.IOS
	puller
	refParser
	resolver
	manifest  manifest.Manifest
	opener    opener
	cachePath string
}
