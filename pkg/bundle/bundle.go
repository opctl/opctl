package bundle

//go:generate counterfeiter -o ./fakeBundle.go --fake-name FakeBundle ./ Bundle

import (
  "github.com/opspec-io/sdk-golang/models"
  "github.com/opspec-io/sdk-golang/util/fs"
  "github.com/opspec-io/sdk-golang/util/format"
)

type Bundle interface {
  CreateCollection(
  req models.CreateCollectionReq,
  ) (err error)

  CreateOp(
  req models.CreateOpReq,
  ) (err error)

  GetCollection(
  collectionBundlePath string,
  ) (
  collectionView models.CollectionView,
  err error,
  )

  GetOp(
  opBundlePath string,
  ) (
  opView models.OpView,
  err error,
  )

  SetCollectionDescription(
  req models.SetCollectionDescriptionReq,
  ) (err error)

  SetOpDescription(
  req models.SetOpDescriptionReq,
  ) (err error)

  TryResolveDefaultCollection(
  req models.TryResolveDefaultCollectionReq,
  ) (
  pathToDefaultCollection string,
  err error,
  )
}

func New() Bundle {
  fileSystem := fs.NewFileSystem()
  yaml := format.NewYamlFormat()
  opViewFactory := newOpViewFactory(fileSystem,yaml)

  return &_bundle{
    collectionViewFactory:newCollectionViewFactory(fileSystem,opViewFactory,yaml),
    fileSystem: fileSystem,
    opViewFactory:opViewFactory,
    yaml: yaml,
  }
}

type _bundle struct {
  collectionViewFactory collectionViewFactory
  fileSystem fs.FileSystem
  opViewFactory opViewFactory
  yaml       format.Format
}
