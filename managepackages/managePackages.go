package managepackages

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ ManagePackages

import (
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/opspec-io/sdk-golang/util/fs"
)

type ManagePackages interface {
	CreatePackage(
		req model.CreatePackageReq,
	) (err error)

	ListPackagesInDir(
		dirPath string,
	) (
		ops []*model.PackageView,
		err error,
	)

	GetPackage(
		packageRef string,
	) (
		packageView model.PackageView,
		err error,
	)

	SetPackageDescription(
		req model.SetPackageDescriptionReq,
	) (err error)
}

func New() ManagePackages {
	fileSystem := fs.NewFileSystem()
	yaml := format.NewYamlFormat()
	packageViewFactory := newPackageViewFactory(fileSystem, yaml)

	return managePackages{
		fileSystem:         fileSystem,
		packageViewFactory: packageViewFactory,
		yaml:               yaml,
	}
}

type managePackages struct {
	fileSystem         fs.FileSystem
	packageViewFactory packageViewFactory
	yaml               format.Format
}
