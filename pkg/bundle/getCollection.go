package bundle

import (
  "github.com/opspec-io/sdk-golang/models"
)

func (this _bundle) GetCollection(
collectionBundlePath string,
) (
collectionView models.CollectionView,
err error,
) {

  return this.collectionViewFactory.Construct(collectionBundlePath)

}
