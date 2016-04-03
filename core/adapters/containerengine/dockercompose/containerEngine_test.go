package dockercompose

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("containerEngine", func() {
  Context(".InitOperation() method", func() {
    It("should invoke compositionRoot.initOperationUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedOperationName := ""

      // wire up fakes
      fakeInitOperationUseCase := new(FakeInitOperationUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.InitOperationUseCaseReturns(fakeInitOperationUseCase)

      objectUnderTest := &_containerEngine{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.InitOperation(providedOperationName)

      /* assert */
      Expect(fakeInitOperationUseCase.ExecuteArgsForCall(0)).To(Equal(providedOperationName))
      Expect(fakeInitOperationUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".RunOperation() method", func() {
    It("should invoke compositionRoot.runOperationUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedOperationName := ""

      fakeRunOperationUseCase := new(FakeRunOperationUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.RunOperationUseCaseReturns(fakeRunOperationUseCase)

      objectUnderTest := &_containerEngine{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.RunOperation(providedOperationName)

      /* assert */
      Expect(fakeRunOperationUseCase.ExecuteArgsForCall(0)).To(Equal(providedOperationName))
      Expect(fakeRunOperationUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
})
