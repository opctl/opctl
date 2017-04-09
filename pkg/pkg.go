// Pkg implements use cases for managing opspec packages
package pkg

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	"github.com/virtual-go/fs"
	"github.com/virtual-go/fs/osfs"
	"github.com/virtual-go/vioutil"
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
	fileSystem := osfs.New()
	ioUtil := vioutil.New(fileSystem)
	validator := newValidator(fileSystem)
	manifestUnmarshaller := newManifestUnmarshaller(ioUtil, validator)

	return pkg{
		fileSystem:           fileSystem,
		getter:               newGetter(fileSystem, manifestUnmarshaller),
		ioUtil:               ioUtil,
		validator:            validator,
		manifestUnmarshaller: manifestUnmarshaller,
	}
}

type pkg struct {
	fileSystem           fs.FS
	getter               getter
	ioUtil               vioutil.VIOUtil
	validator            validator
	manifestUnmarshaller manifestUnmarshaller
}
