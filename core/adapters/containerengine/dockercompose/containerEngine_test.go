package dockercompose

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("containerEngine", func() {
  Context(".InitOp() method", func() {
    It("should invoke compositionRoot.initOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedPathToOpDir := ""
      providedOpName := ""

      // wire up fakes
      fakeInitOpUseCase := new(fakeInitOpUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.InitOpUseCaseReturns(fakeInitOpUseCase)

      objectUnderTest := &_containerEngine{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.InitOp(
        providedPathToOpDir,
        providedOpName,
      )

      /* assert */
      Expect(fakeInitOpUseCase.ExecuteArgsForCall(0)).To(Equal(providedOpName))
      Expect(fakeInitOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".RunOp() method", func() {
    It("should invoke compositionRoot.runOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedPathToOpDir := ""
      providedOpName := ""

      fakeRunOpUseCase := new(fakeRunOpUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.RunOpUseCaseReturns(fakeRunOpUseCase)

      objectUnderTest := &_containerEngine{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.RunOp(
        providedPathToOpDir,
        providedOpName,
      )

      /* assert */
      Expect(fakeRunOpUseCase.ExecuteArgsForCall(0)).To(Equal(providedOpName))
      Expect(fakeRunOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
})
