package core

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/dev-op-spec/engine/core/models"
)

var _ = Describe("_sdk", func() {
  Context(".AddOperation() method", func() {
    It("should invoke compositionRoot.addOperationUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedAddOperationReq := models.NewAddOperationReq(&models.ProjectUrl{}, "", "")

      // wire up fakes
      fakeAddOperationUseCase := new(fakeAddOperationUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.AddOperationUseCaseReturns(fakeAddOperationUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.AddOperation(*providedAddOperationReq)

      /* assert */
      Expect(fakeAddOperationUseCase.ExecuteArgsForCall(0)).To(Equal(*providedAddOperationReq))
      Expect(fakeAddOperationUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".AddSubOperation() method", func() {
    It("should invoke compositionRoot.addSubOperationUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedAddSubOperationReq := models.NewAddSubOperationReq(&models.ProjectUrl{}, false, "", "", "")

      // wire up fakes
      fakeAddSubOperationUseCase := new(fakeAddSubOperationUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.AddSubOperationUseCaseReturns(fakeAddSubOperationUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.AddSubOperation(*providedAddSubOperationReq)

      /* assert */
      Expect(fakeAddSubOperationUseCase.ExecuteArgsForCall(0)).To(Equal(*providedAddSubOperationReq))
      Expect(fakeAddSubOperationUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".ListOperations() method", func() {
    It("should invoke compositionRoot.listOperationsUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedProjectUrl := &models.ProjectUrl{}
      expectedReturnedOperations := make([]models.OperationView, 0)

      // wire up fakes
      fakeListOperationsUseCase := new(fakeListOperationsUseCase)
      fakeListOperationsUseCase.ExecuteReturns(expectedReturnedOperations, nil)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.ListOperationsUseCaseReturns(fakeListOperationsUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      actualReturnedOperations, _ := objectUnderTest.ListOperations(providedProjectUrl)

      /* assert */
      Expect(actualReturnedOperations).To(Equal(expectedReturnedOperations))

    })
  })
  Context(".RunOperation() method", func() {
    It("should invoke compositionRoot.runOperationUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedRunOperationReq := models.NewRunOperationReq(&models.ProjectUrl{}, "")

      // wire up fakes
      fakeRunOperationUseCase := new(fakeRunOperationUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.RunOperationUseCaseReturns(fakeRunOperationUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.RunOperation(*providedRunOperationReq)

      /* assert */
      executeArg0, executeArg1 := fakeRunOperationUseCase.ExecuteArgsForCall(0)
      Expect(executeArg0).To(Equal(*providedRunOperationReq))
      Expect(executeArg1).To(Equal(make([]string, 0)))
      Expect(fakeRunOperationUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".SetDescriptionOfOperation() method", func() {
    It("should invoke compositionRoot.setDescriptionOfOperationUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedSetDescriptionOfOperationReq := models.NewSetDescriptionOfOperationReq(&models.ProjectUrl{}, "", "")

      // wire up fakes
      fakeSetDescriptionOfOperationUseCase := new(fakeSetDescriptionOfOperationUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.SetDescriptionOfOperationUseCaseReturns(fakeSetDescriptionOfOperationUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.SetDescriptionOfOperation(*providedSetDescriptionOfOperationReq)

      /* assert */
      Expect(fakeSetDescriptionOfOperationUseCase.ExecuteArgsForCall(0)).To(Equal(*providedSetDescriptionOfOperationReq))
      Expect(fakeSetDescriptionOfOperationUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })

})
