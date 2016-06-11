package opspec

//go:generate counterfeiter -o ./fakeSdk.go --fake-name FakeSdk ./ Sdk

import "github.com/opspec-io/sdk-golang/models"

type Sdk interface {
  CreateOp(
  req models.CreateOpReq,
  ) (err error)

  SetCollectionDescription(
  req models.SetCollectionDescriptionReq,
  ) (err error)

  SetOpDescription(
  req models.SetOpDescriptionReq,
  ) (err error)
}

func New(
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

func (this _sdk) CreateOp(
req models.CreateOpReq,
) (err error) {
  return this.
  compositionRoot.
  CreateOpUseCase().
  Execute(req)
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
