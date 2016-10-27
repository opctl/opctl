package opspec

//go:generate counterfeiter -o ./fakeSdk.go --fake-name FakeSdk ./ Sdk

import "github.com/opspec-io/sdk-golang/models"

type Sdk interface {
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

func New(
) Sdk {
  return newSdk(
    newFilesystem(),
  )
}

func newSdk(
filesystem Filesystem,
) (sdk Sdk) {

  var compositionRoot compositionRoot
  compositionRoot = newCompositionRoot(
    filesystem,
  )

  sdk = &_sdk{
    compositionRoot:compositionRoot,
  }

  return
}

type _sdk struct {
  compositionRoot compositionRoot
}

func (this _sdk) CreateCollection(
req models.CreateCollectionReq,
) (err error) {
  return this.
  compositionRoot.
  CreateCollectionUseCase().
  Execute(req)
}

func (this _sdk) CreateOp(
req models.CreateOpReq,
) (err error) {
  return this.
  compositionRoot.
  CreateOpUseCase().
  Execute(req)
}

func (this _sdk) GetCollection(
collectionBundlePath string,
) (
collectionView models.CollectionView,
err error,
) {
  return this.
  compositionRoot.
  GetCollectionUseCase().
  Execute(collectionBundlePath)
}

func (this _sdk) GetOp(
opBundlePath string,
) (
opView models.OpView,
err error,
) {
  return this.
  compositionRoot.
  GetOpUseCase().
  Execute(opBundlePath)
}

func (this _sdk) SetCollectionDescription(
req models.SetCollectionDescriptionReq,
) (err error) {
  return this.
  compositionRoot.
  SetCollectionDescriptionUseCase().
  Execute(req)
}

func (this _sdk) SetOpDescription(
req models.SetOpDescriptionReq,
) (err error) {
  return this.
  compositionRoot.
  SetOpDescriptionUseCase().
  Execute(req)
}

func (this _sdk) TryResolveDefaultCollection(
req models.TryResolveDefaultCollectionReq,
) (
pathToDefaultCollection string,
err error,
) {
  return this.
  compositionRoot.
  TryResolveDefaultCollectionUseCase().
  Execute(req)
}
