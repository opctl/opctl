package dockercompose

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("containerEngine", func() {
  Context(".InitDevOp() method", func() {
    It("should invoke compositionRoot.initDevOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedDevOpName := ""

      // wire up fakes
      fakeInitDevOpUseCase := new(FakeInitDevOpUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.InitDevOpUseCaseReturns(fakeInitDevOpUseCase)

      objectUnderTest := &_containerEngine{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.InitDevOp(providedDevOpName)

      /* assert */
      Expect(fakeInitDevOpUseCase.ExecuteArgsForCall(0)).To(Equal(providedDevOpName))
      Expect(fakeInitDevOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".RunDevOp() method", func() {
    It("should invoke compositionRoot.runDevOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedDevOpName := ""

      fakeRunDevOpUseCase := new(FakeRunDevOpUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.RunDevOpUseCaseReturns(fakeRunDevOpUseCase)

      objectUnderTest := &_containerEngine{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.RunDevOp(providedDevOpName)

      /* assert */
      Expect(fakeRunDevOpUseCase.ExecuteArgsForCall(0)).To(Equal(providedDevOpName))
      Expect(fakeRunDevOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
})
