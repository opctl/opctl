// Pkg implements use cases for managing opspec packages
package pkg

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	fsPkg "github.com/virtual-go/fs"
	"github.com/virtual-go/fs/osfs"
	"github.com/virtual-go/gopkg.in-src-d-go-git.v4"
	"github.com/virtual-go/vioutil"
	"github.com/virtual-go/vos"
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
		req *PullOpts,
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
	fs := osfs.New()
	ioUtil := vioutil.New(fs)
	manifestValidator := newManifestValidator(fs)
	manifestUnmarshaller := newManifestUnmarshaller(ioUtil, manifestValidator)

	return pkg{
		fs:                   fs,
		git:                  vgit.New(),
		ioUtil:               ioUtil,
		os:                   vos.New(fs),
		manifestUnmarshaller: manifestUnmarshaller,
		manifestValidator:    manifestValidator,
	}
}

type pkg struct {
	fs                   fsPkg.FS
	git                  vgit.VGit
	ioUtil               vioutil.VIOUtil
	os                   vos.VOS
	manifestValidator    manifestValidator
	manifestUnmarshaller manifestUnmarshaller
}
