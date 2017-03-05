package pkg

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
)

func (this pkg) GetCollection(
	collectionPackagePath string,
) (
	collectionView model.CollectionView,
	err error,
) {

	return this.collectionViewFactory.Construct(collectionPackagePath)

}
