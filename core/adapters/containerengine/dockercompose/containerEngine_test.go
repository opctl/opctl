package dockercompose

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/dev-op-spec/engine/core/models"
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
      providedLogChannel := make(chan *models.LogEntry)

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
        providedLogChannel,
      )

      /* assert */
      receivedPathToOpDir, receivedOpName, receivedLogChannel := fakeRunOpUseCase.ExecuteArgsForCall(0)
      Expect(receivedPathToOpDir).To(Equal(providedPathToOpDir))
      Expect(receivedOpName).To(Equal(providedOpName))
      Expect(receivedLogChannel).To(Equal(providedLogChannel))
      Expect(fakeRunOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
})
