package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Describe("_getCollection", func() {

	Context("Execute", func() {

		It("should call collectionViewFactory.Construct with expected args", func() {

			/* arrange */

			providedCollectionPackagePath := "/dummy/path"

			fakeCollectionViewFactory := new(fakeCollectionViewFactory)

			objectUnderTest := &pkg{
				collectionViewFactory: fakeCollectionViewFactory,
			}

			/* act */
			objectUnderTest.GetCollection(
				providedCollectionPackagePath,
			)

			/* assert */
			Expect(fakeCollectionViewFactory.ConstructArgsForCall(0)).To(Equal(providedCollectionPackagePath))

		})

		It("should return result of collectionViewFactory.Construct", func() {

			/* arrange */
			expectedCollectionView := model.CollectionView{
				Description: "dummyDescription",
				Name:        "dummyName",
				Ops:         []model.OpView{},
			}
			expectedError := errors.New("ConstructError")

			fakeCollectionViewFactory := new(fakeCollectionViewFactory)
			fakeCollectionViewFactory.ConstructReturns(expectedCollectionView, expectedError)

			objectUnderTest := &pkg{
				collectionViewFactory: fakeCollectionViewFactory,
			}

			/* act */
			actualCollectionView, actualError := objectUnderTest.GetCollection("/dummy/path")

			/* assert */
			Expect(actualCollectionView).To(Equal(expectedCollectionView))
			Expect(actualError).To(Equal(expectedError))

		})

	})

})
