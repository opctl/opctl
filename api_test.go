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
  Context(".AddOp() method", func() {
    It("should invoke compositionRoot.addOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedAddOpReq := models.NewAddOpReq("", "", "")

      // wire up fakes
      fakeAddOpUseCase := new(fakeAddOpUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.AddOpUseCaseReturns(fakeAddOpUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.AddOp(*providedAddOpReq)

      /* assert */
      Expect(fakeAddOpUseCase.ExecuteArgsForCall(0)).To(Equal(*providedAddOpReq))
      Expect(fakeAddOpUseCase.ExecuteCallCount()).To(Equal(1))

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
