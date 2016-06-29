package opspec

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/models"
  "errors"
)

var _ = Describe("_getOpUseCase", func() {

  Context("Execute", func() {

    It("should call opViewFactory.Construct with expected args", func() {

      /* arrange */

      providedOpBundlePath := "/dummy/path"

      fakeOpViewFactory := new(fakeOpViewFactory)

      objectUnderTest := newGetOpUseCase(fakeOpViewFactory)

      /* act */
      objectUnderTest.Execute(
        providedOpBundlePath,
      )

      /* assert */
      Expect(fakeOpViewFactory.ConstructArgsForCall(0)).To(Equal(providedOpBundlePath))

    })

    It("should return result of opViewFactory.Construct", func() {

      /* arrange */
      expectedOpView := *models.NewOpView(
        "dummy description",
        "dummy name",
        []models.OpParamView{},
        []models.SubOpView{},
      )
      expectedError := errors.New("ConstructError")

      fakeOpViewFactory := new(fakeOpViewFactory)
      fakeOpViewFactory.ConstructReturns(expectedOpView, expectedError)

      objectUnderTest := newGetOpUseCase(fakeOpViewFactory)

      /* act */
      actualOpView, actualError := objectUnderTest.Execute("/dummy/path")

      /* assert */
      Expect(actualOpView).To(Equal(expectedOpView))
      Expect(actualError).To(Equal(expectedError))

    })

  })

})
