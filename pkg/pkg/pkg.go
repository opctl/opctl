package pkg

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Package

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/opspec-io/sdk-golang/util/fs"
)

type Pkg interface {
	CreateCollection(
		req model.CreateCollectionReq,
	) (err error)

	CreateOp(
		req model.CreateOpReq,
	) (err error)

	GetCollection(
		collectionPackagePath string,
	) (
		collectionView model.CollectionView,
		err error,
	)

	GetOp(
		opPackagePath string,
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

func New() Pkg {
	fileSystem := fs.NewFileSystem()
	yaml := format.NewYamlFormat()
	opViewFactory := newOpViewFactory(fileSystem, yaml)

	return &pkg{
		collectionViewFactory: newCollectionViewFactory(fileSystem, opViewFactory, yaml),
		fileSystem:            fileSystem,
		opViewFactory:         opViewFactory,
		yaml:                  yaml,
	}
}

type pkg struct {
	collectionViewFactory collectionViewFactory
	fileSystem            fs.FileSystem
	opViewFactory         opViewFactory
	yaml                  format.Format
}
