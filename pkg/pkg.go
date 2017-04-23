// Pkg implements use cases for managing opspec packages
package pkg

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/vgit"
	fsPkg "github.com/virtual-go/fs"
	"github.com/virtual-go/fs/osfs"
	"github.com/virtual-go/vioutil"
	"github.com/virtual-go/vos"
)

type Pkg interface {
	// Create creates an opspec package
	Create(
		req CreateReq,
	) error

	// Resolve resolves a local package according to opspec package resolution rules and returns it's absolute path.
	Resolve(
		basePath,
		pkgRef string,
	) (string, bool)

	// Pull pulls a package from a remote source
	Pull(
		pkgRef string,
		req *PullOpts,
	) error

	// Get gets a local package
	Get(
		basePath,
		pkgRef string,
	) (*model.PkgManifest, error)

	// List lists packages according to opspec package resolution rules
	List(
		dirPath string,
	) ([]*model.PkgManifest, error)

	// SetDescription sets the description of a package
	SetDescription(
		req SetDescriptionReq,
	) error

	// Validate validates an opspec package
	Validate(
		pkgRef string,
	) []error
}

func New() Pkg {
	fs := osfs.New()
	os := vos.New(fs)
	ioUtil := vioutil.New(fs)
	validator := newValidator(fs)
	resolver := newResolver(os)
	manifestUnmarshaller := newManifestUnmarshaller(ioUtil, validator)

	return pkg{
		fs:                   fs,
		getter:               newGetter(manifestUnmarshaller, resolver),
		git:                  vgit.New(),
		ioUtil:               ioUtil,
		manifestUnmarshaller: manifestUnmarshaller,
		resolver:             resolver,
		validator:            validator,
	}
}

type pkg struct {
	fs                   fsPkg.FS
	getter               getter
	git                  vgit.VGit
	ioUtil               vioutil.VIOUtil
	resolver             resolver
	validator            validator
	manifestUnmarshaller manifestUnmarshaller
}
