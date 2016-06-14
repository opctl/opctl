package opspec

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/models"
  "errors"
)

var _ = Describe("_getCollectionUseCase", func() {

  Context("Execute", func() {

    It("should call collectionViewFactory.Construct with expected args", func() {

      /* arrange */

      providedCollectionBundlePath := "/dummy/path"

      fakeCollectionViewFactory := new(fakeCollectionViewFactory)

      objectUnderTest := newGetCollectionUseCase(fakeCollectionViewFactory)

      /* act */
      objectUnderTest.Execute(
        providedCollectionBundlePath,
      )

      /* assert */
      Expect(fakeCollectionViewFactory.ConstructArgsForCall(0)).To(Equal(providedCollectionBundlePath))

    })

    It("should return result of collectionViewFactory.Construct", func() {

      /* arrange */
      expectedCollectionView := *models.NewCollectionView(
        "dummy description",
        "dummy name",
        []models.OpView{},
      )
      expectedError := errors.New("ConstructError")

      fakeCollectionViewFactory := new(fakeCollectionViewFactory)
      fakeCollectionViewFactory.ConstructReturns(expectedCollectionView, expectedError)

      objectUnderTest := newGetCollectionUseCase(fakeCollectionViewFactory)

      /* act */
      actualCollectionView, actualError := objectUnderTest.Execute("/dummy/path")

      /* assert */
      Expect(actualCollectionView).To(Equal(expectedCollectionView))
      Expect(actualError).To(Equal(expectedError))

    })

  })

})
