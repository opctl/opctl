package pkg

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/opspec-io/sdk-golang/util/fs"
	"github.com/xeipuuv/gojsonschema"
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

	manifestSchemaBytes, err := pkgDataPackageManifestSchemaJsonBytes()
	if nil != err {
		panic(err)
	}

	manifestSchema, err := gojsonschema.NewSchema(
		gojsonschema.NewBytesLoader(manifestSchemaBytes),
	)
	if err != nil {
		panic(err)
	}

	return pkg{
		fileSystem:         fileSystem,
		manifestSchema:     manifestSchema,
		packageViewFactory: packageViewFactory,
		yaml:               yaml,
	}
}

type pkg struct {
	fileSystem         fs.FileSystem
	manifestSchema     *gojsonschema.Schema
	packageViewFactory packageViewFactory
	yaml               format.Format
}
