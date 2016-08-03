package opspec

//go:generate counterfeiter -o ./fakeGetCollectionUseCase.go --fake-name fakeGetCollectionUseCase ./ getCollectionUseCase

import (
  "github.com/opspec-io/sdk-golang/models"
)

type getCollectionUseCase interface {
  Execute(
  collectionBundlePath string,
  ) (
  collectionView models.CollectionView,
  err error,
  )
}

func newGetCollectionUseCase(
collectionViewFactory collectionViewFactory,
) getCollectionUseCase {

  return &_getCollectionUseCase{
    collectionViewFactory:collectionViewFactory,
  }

}

type _getCollectionUseCase struct {
  collectionViewFactory collectionViewFactory
}

func (this _getCollectionUseCase) Execute(
collectionBundlePath string,
) (
collectionView models.CollectionView,
err error,
) {

  return this.collectionViewFactory.Construct(collectionBundlePath)

}
