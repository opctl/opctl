package pkg

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/opspec-io/sdk-golang/util/fs"
)

type Pkg interface {
	Create(
		req CreateReq,
	) (err error)

	Get(
		pkgRef string,
	) (
		packageView model.PackageView,
		err error,
	)

	ListPackagesInDir(
		dirPath string,
	) (
		ops []*model.PackageView,
		err error,
	)

	SetDescription(
		req SetDescriptionReq,
	) (err error)

	Validate(
		pkgRef string,
	) (errs []error)
}

func New() Pkg {
	fileSystem := fs.NewFileSystem()
	yaml := format.NewYamlFormat()
	packageViewFactory := newPackageViewFactory(fileSystem, yaml)

	return pkg{
		fileSystem:         fileSystem,
		packageViewFactory: packageViewFactory,
		yaml:               yaml,
		validator:          newValidator(),
	}
}

type pkg struct {
	fileSystem         fs.FileSystem
	packageViewFactory packageViewFactory
	yaml               format.Format
	validator          validator
}
