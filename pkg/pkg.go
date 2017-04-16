// Pkg implements use cases for managing opspec packages
package pkg

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Pkg

import (
	"github.com/opspec-io/sdk-golang/model"
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

	// Get gets a package according to opspec package resolution rules
	Get(
		req *GetReq,
	) (*model.PkgManifest, error)

	// ResolveRef resolves a pkgRef according to opspec package resolution rules.
	ResolveRef(
		pkgRef string,
	) string

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
	refResolver := newRefResolver(os)
	manifestUnmarshaller := newManifestUnmarshaller(ioUtil, validator)

	return pkg{
		fs:                   fs,
		getter:               newGetter(fs, manifestUnmarshaller, refResolver),
		ioUtil:               ioUtil,
		manifestUnmarshaller: manifestUnmarshaller,
		refResolver:          refResolver,
		validator:            validator,
	}
}

type pkg struct {
	fs                   fsPkg.FS
	getter               getter
	ioUtil               vioutil.VIOUtil
	refResolver          refResolver
	validator            validator
	manifestUnmarshaller manifestUnmarshaller
}
