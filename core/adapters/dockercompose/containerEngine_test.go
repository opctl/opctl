package dockercompose

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("containerEngine", func() {
  Context(".InitDevOp() method", func() {
    It("should invoke compositionRoot.initDevOpUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */
      providedDevOpName := ""

      // wire up fakes
      fakeInitDevOpUCExecuter := new(FakeInitDevOpUcExecuter)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.InitDevOpUcExecuterReturns(fakeInitDevOpUCExecuter)

      objectUnderTest := &containerEngineImpl{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.InitDevOp(providedDevOpName)

      /* assert */
      Expect(fakeInitDevOpUCExecuter.ExecuteArgsForCall(0)).To(Equal(providedDevOpName))
      Expect(fakeInitDevOpUCExecuter.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".RunDevOp() method", func() {
    It("should invoke compositionRoot.runDevOpUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */
      providedDevOpName := ""

      fakeRunDevOpUcExecuter := new(FakeRunDevOpUcExecuter)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.RunDevOpUcExecuterReturns(fakeRunDevOpUcExecuter)

      objectUnderTest := &containerEngineImpl{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.RunDevOp(providedDevOpName)

      /* assert */
      Expect(fakeRunDevOpUcExecuter.ExecuteArgsForCall(0)).To(Equal(providedDevOpName))
      Expect(fakeRunDevOpUcExecuter.ExecuteCallCount()).To(Equal(1))

    })
  })
})
