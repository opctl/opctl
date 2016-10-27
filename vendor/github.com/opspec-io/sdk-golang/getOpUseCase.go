package opspec

//go:generate counterfeiter -o ./fakeGetOpUseCase.go --fake-name fakeGetOpUseCase ./ getOpUseCase

import (
  "github.com/opspec-io/sdk-golang/models"
)

type getOpUseCase interface {
  Execute(
  opBundlePath string,
  ) (
  opView models.OpView,
  err error,
  )
}

func newGetOpUseCase(
opViewFactory opViewFactory,
) getOpUseCase {

  return &_getOpUseCase{
    opViewFactory:opViewFactory,
  }

}

type _getOpUseCase struct {
  opViewFactory opViewFactory
}

func (this _getOpUseCase) Execute(
opBundlePath string,
) (
opView models.OpView,
err error,
) {

  return this.opViewFactory.Construct(opBundlePath)

}
