package sdk

//go:generate counterfeiter -o ./fakeApi.go --fake-name FakeApi ./ Api

import "github.com/opspec-io/sdk-golang/models"

type Api interface {
  SetDescriptionOfOp(
  req models.SetDescriptionOfOpReq,
  ) (err error)
}

func New(
filesys Filesystem,
) (api Api, err error) {

  var compositionRoot compositionRoot
  compositionRoot, err = newCompositionRoot(
    filesys,
  )
  if (nil != err) {
    return
  }

  api = &_api{
    compositionRoot:compositionRoot,
  }

  return
}

type _api struct {
  compositionRoot compositionRoot
}

func (this _api) SetDescriptionOfOp(
req models.SetDescriptionOfOpReq,
) (err error) {
  return this.
  compositionRoot.
  SetDescriptionOfOpUseCase().
  Execute(req)
}
