package bundle

import (
  "github.com/opspec-io/sdk-golang/pkg/model"
)

func (this _bundle) GetCollection(
collectionBundlePath string,
) (
collectionView model.CollectionView,
err error,
) {

  return this.collectionViewFactory.Construct(collectionBundlePath)

}
