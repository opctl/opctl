package sdk

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/models"
)

var _ = Describe("_api", func() {

  var fakeFilesystem = new(FakeFilesystem)

  Context("new()", func() {
    It("should return an instance of _api", func() {

      /* arrange/act */
      objectUnderTest := New(
        fakeFilesystem,
      )

      /* assert */
      Expect(objectUnderTest).To(BeAssignableToTypeOf(&_api{}))

    })
  })
  Context(".CreateOp() method", func() {
    It("should invoke compositionRoot.createOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedCreateOpReq := models.NewCreateOpReq("", "", "")

      // wire up fakes
      fakeCreateOpUseCase := new(fakeCreateOpUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.CreateOpUseCaseReturns(fakeCreateOpUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.CreateOp(*providedCreateOpReq)

      /* assert */
      Expect(fakeCreateOpUseCase.ExecuteArgsForCall(0)).To(Equal(*providedCreateOpReq))
      Expect(fakeCreateOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".SetDescriptionOfOp() method", func() {
    It("should invoke compositionRoot.setDescriptionOfOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedSetDescriptionOfOpReq := models.NewSetDescriptionOfOpReq("", "")

      // wire up fakes
      fakeSetDescriptionOfOpUseCase := new(fakeSetDescriptionOfOpUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.SetDescriptionOfOpUseCaseReturns(fakeSetDescriptionOfOpUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.SetDescriptionOfOp(*providedSetDescriptionOfOpReq)

      /* assert */
      Expect(fakeSetDescriptionOfOpUseCase.ExecuteArgsForCall(0)).To(Equal(*providedSetDescriptionOfOpReq))
      Expect(fakeSetDescriptionOfOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })

})
