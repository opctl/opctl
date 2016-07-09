package dockercompose

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opctl/engine/core/logging"
)

var _ = Describe("containerEngine", func() {
  Context(".RunOp() method", func() {
    It("should invoke compositionRoot.runOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedOpArgs := map[string]string{}
      providedCorrelationId := ""
      providedOpBundlePath := ""
      providedOpName := ""
      providedOpRunId := ""
      providedLogger := new(logging.Logger)

      fakeRunOpUseCase := new(fakeRunOpUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.RunOpUseCaseReturns(fakeRunOpUseCase)

      objectUnderTest := &_containerEngine{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.RunOp(
        providedCorrelationId,
        providedOpArgs,
        providedOpBundlePath,
        providedOpName,
        providedOpRunId,
        *providedLogger,
      )

      /* assert */
      receivedCorrelationId,
      receivedArgs,
      receivedOpBundlePath,
      receivedOpName,
      receivedOpRunId,
      receivedLogger := fakeRunOpUseCase.ExecuteArgsForCall(0)
      Expect(receivedArgs).To(Equal(providedOpArgs))
      Expect(receivedCorrelationId).To(Equal(providedCorrelationId))
      Expect(receivedOpBundlePath).To(Equal(providedOpBundlePath))
      Expect(receivedOpName).To(Equal(providedOpName))
      Expect(receivedOpRunId).To(Equal(providedOpRunId))
      Expect(receivedLogger).To(Equal(*providedLogger))
      Expect(fakeRunOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
})
