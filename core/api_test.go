package core

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/dev-op-spec/engine/core/models"
)

var _ = Describe("_sdk", func() {
  Context(".AddOp() method", func() {
    It("should invoke compositionRoot.addOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedAddOpReq := models.NewAddOpReq(&models.Url{}, "", "")

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
  Context(".AddSubOp() method", func() {
    It("should invoke compositionRoot.addSubOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedAddSubOpReq := models.NewAddSubOpReq(&models.Url{}, "", "", "")

      // wire up fakes
      fakeAddSubOpUseCase := new(fakeAddSubOpUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.AddSubOpUseCaseReturns(fakeAddSubOpUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.AddSubOp(*providedAddSubOpReq)

      /* assert */
      Expect(fakeAddSubOpUseCase.ExecuteArgsForCall(0)).To(Equal(*providedAddSubOpReq))
      Expect(fakeAddSubOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".ListOps() method", func() {
    It("should invoke compositionRoot.listOpsUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedProjectUrl := &models.Url{}
      expectedReturnedOps := make([]models.OpDetailedView, 0)

      // wire up fakes
      fakeListOpsUseCase := new(fakeListOpsUseCase)
      fakeListOpsUseCase.ExecuteReturns(expectedReturnedOps, nil)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.ListOpsUseCaseReturns(fakeListOpsUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      actualReturnedOps, _ := objectUnderTest.ListOps(providedProjectUrl)

      /* assert */
      Expect(actualReturnedOps).To(Equal(expectedReturnedOps))

    })
  })
  Context(".RunOp() method", func() {
    It("should invoke compositionRoot.runOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedRunOpReq := models.NewRunOpReq(&models.Url{})

      // wire up fakes
      fakeRunOpUseCase := new(fakeRunOpUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.RunOpUseCaseReturns(fakeRunOpUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.RunOp(*providedRunOpReq)

      /* assert */
      executeArg0, executeArg1 := fakeRunOpUseCase.ExecuteArgsForCall(0)
      Expect(executeArg0).To(Equal(*providedRunOpReq))
      Expect(executeArg1).To(Equal(make([]*models.Url, 0)))
      Expect(fakeRunOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".SetDescriptionOfOp() method", func() {
    It("should invoke compositionRoot.setDescriptionOfOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedSetDescriptionOfOpReq := models.NewSetDescriptionOfOpReq(&models.Url{}, new(string), "")

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
