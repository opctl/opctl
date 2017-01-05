package bundle

//go:generate counterfeiter -o ./fakeBundle.go --fake-name FakeBundle ./ Bundle

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/opspec-io/sdk-golang/util/fs"
)

type Bundle interface {
	CreateCollection(
		req model.CreateCollectionReq,
	) (err error)

	CreateOp(
		req model.CreateOpReq,
	) (err error)

	GetCollection(
		collectionBundlePath string,
	) (
		collectionView model.CollectionView,
		err error,
	)

	GetOp(
		opBundlePath string,
	) (
		opView model.OpView,
		err error,
	)

	SetCollectionDescription(
		req model.SetCollectionDescriptionReq,
	) (err error)

	SetOpDescription(
		req model.SetOpDescriptionReq,
	) (err error)

	TryResolveDefaultCollection(
		req model.TryResolveDefaultCollectionReq,
	) (
		pathToDefaultCollection string,
		err error,
	)
}

func New() Bundle {
	fileSystem := fs.NewFileSystem()
	yaml := format.NewYamlFormat()
	opViewFactory := newOpViewFactory(fileSystem, yaml)

	return &_bundle{
		collectionViewFactory: newCollectionViewFactory(fileSystem, opViewFactory, yaml),
		fileSystem:            fileSystem,
		opViewFactory:         opViewFactory,
		yaml:                  yaml,
	}
}

type _bundle struct {
	collectionViewFactory collectionViewFactory
	fileSystem            fs.FileSystem
	opViewFactory         opViewFactory
	yaml                  format.Format
}
