package pkg

//go:generate counterfeiter -o ./fakeOpViewFactory.go --fake-name fakeOpViewFactory ./ opViewFactory

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/opspec-io/sdk-golang/util/fs"
	"path"
)

type opViewFactory interface {
	Construct(
		opPackagePath string,
	) (
		opView model.OpView,
		err error,
	)
}

func newOpViewFactory(
	fileSystem fs.FileSystem,
	yaml format.Format,
) opViewFactory {

	return &_opViewFactory{
		fileSystem: fileSystem,
		yaml:       yaml,
	}

}

type _opViewFactory struct {
	fileSystem fs.FileSystem
	yaml       format.Format
}

func (this _opViewFactory) Construct(
	opPackagePath string,
) (
	opView model.OpView,
	err error,
) {

	opManifestPath := path.Join(opPackagePath, NameOfOpManifestFile)

	opManifestBytes, err := this.fileSystem.GetBytesOfFile(
		opManifestPath,
	)
	if nil != err {
		return
	}

	opManifest := model.OpManifest{}
	err = this.yaml.To(
		opManifestBytes,
		&opManifest,
	)
	if nil != err {
		return
	}

	opView = model.OpView{
		Description: opManifest.Description,
		Inputs:      opManifest.Inputs,
		Name:        opManifest.Name,
		Outputs:     opManifest.Outputs,
		Run:         opManifest.Run,
		Version:     opManifest.Version,
	}

	return

}
