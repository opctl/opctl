package bundle

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/pkg/model"
  "errors"
)

var _ = Describe("_getCollection", func() {

  Context("Execute", func() {

    It("should call collectionViewFactory.Construct with expected args", func() {

      /* arrange */

      providedCollectionBundlePath := "/dummy/path"

      fakeCollectionViewFactory := new(fakeCollectionViewFactory)

      objectUnderTest := &_bundle{
        collectionViewFactory:fakeCollectionViewFactory,
      }

      /* act */
      objectUnderTest.GetCollection(
        providedCollectionBundlePath,
      )

      /* assert */
      Expect(fakeCollectionViewFactory.ConstructArgsForCall(0)).To(Equal(providedCollectionBundlePath))

    })

    It("should return result of collectionViewFactory.Construct", func() {

      /* arrange */
      expectedCollectionView := *model.NewCollectionView(
        "dummy description",
        "dummy name",
        []model.OpView{},
      )
      expectedError := errors.New("ConstructError")

      fakeCollectionViewFactory := new(fakeCollectionViewFactory)
      fakeCollectionViewFactory.ConstructReturns(expectedCollectionView, expectedError)

      objectUnderTest := &_bundle{
        collectionViewFactory:fakeCollectionViewFactory,
      }

      /* act */
      actualCollectionView, actualError := objectUnderTest.GetCollection("/dummy/path")

      /* assert */
      Expect(actualCollectionView).To(Equal(expectedCollectionView))
      Expect(actualError).To(Equal(expectedError))

    })

  })

})
