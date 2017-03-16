package managepackages

//go:generate counterfeiter -o ./fakePackageViewFactory.go --fake-name fakePackageViewFactory ./ packageViewFactory

import (
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/opspec-io/sdk-golang/util/fs"
	"path"
)

type packageViewFactory interface {
	Construct(
		packageRef string,
	) (
		packageView model.PackageView,
		err error,
	)
}

func newPackageViewFactory(
	fileSystem fs.FileSystem,
	yaml format.Format,
) packageViewFactory {

	return &_packageViewFactory{
		fileSystem: fileSystem,
		yaml:       yaml,
	}

}

type _packageViewFactory struct {
	fileSystem fs.FileSystem
	yaml       format.Format
}

func (this _packageViewFactory) Construct(
	packageRef string,
) (
	packageView model.PackageView,
	err error,
) {

	packageManifestViewPath := path.Join(packageRef, NameOfPackageManifestFile)

	packageManifestViewBytes, err := this.fileSystem.GetBytesOfFile(
		packageManifestViewPath,
	)
	if nil != err {
		return
	}

	packageManifestView := model.PackageManifestView{}
	err = this.yaml.To(
		packageManifestViewBytes,
		&packageManifestView,
	)
	if nil != err {
		return
	}

	packageView = model.PackageView{
		Description: packageManifestView.Description,
		Inputs:      packageManifestView.Inputs,
		Name:        packageManifestView.Name,
		Outputs:     packageManifestView.Outputs,
		Run:         packageManifestView.Run,
		Version:     packageManifestView.Version,
	}

	return

}
