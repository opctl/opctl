package dockercompose

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opctl/engine/core/logging"
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
      providedArgs := map[string]string{}
      providedCorrelationId := ""
      providedPathToOpDir := ""
      providedOpName := ""
      providedLogger := new(logging.Logger)

      fakeRunOpUseCase := new(fakeRunOpUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.RunOpUseCaseReturns(fakeRunOpUseCase)

      objectUnderTest := &_containerEngine{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.RunOp(
        providedArgs,
        providedCorrelationId,
        providedPathToOpDir,
        providedOpName,
        *providedLogger,
      )

      /* assert */
      receivedArgs,
      receivedCorrelationId,
      receivedPathToOpDir,
      receivedOpName,
      receivedLogger := fakeRunOpUseCase.ExecuteArgsForCall(0)
      Expect(receivedArgs).To(Equal(providedArgs))
      Expect(receivedCorrelationId).To(Equal(providedCorrelationId))
      Expect(receivedPathToOpDir).To(Equal(providedPathToOpDir))
      Expect(receivedOpName).To(Equal(providedOpName))
      Expect(receivedLogger).To(Equal(*providedLogger))
      Expect(fakeRunOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
})
