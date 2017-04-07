package pkg

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/virtual-go/fs"
	"github.com/virtual-go/fs/osfs"
	"github.com/virtual-go/vioutil"
)

type Pkg interface {
	Create(
		req CreateReq,
	) (err error)

	Get(
		req *GetReq,
	) (
		packageView *model.PackageView,
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
	fileSystem := osfs.New()
	ioUtil := vioutil.New(fileSystem)
	validator := newValidator(fileSystem)
	yaml := format.NewYamlFormat()
	viewFactory := newViewFactory(ioUtil, validator, yaml)

	return pkg{
		fileSystem:  fileSystem,
		ioUtil:      ioUtil,
		getter:      newGetter(fileSystem, viewFactory),
		viewFactory: viewFactory,
		yaml:        yaml,
		validator:   validator,
	}
}

type pkg struct {
	fileSystem  fs.FS
	ioUtil      vioutil.VIOUtil
	getter      getter
	viewFactory viewFactory
	yaml        format.Format
	validator   validator
}
