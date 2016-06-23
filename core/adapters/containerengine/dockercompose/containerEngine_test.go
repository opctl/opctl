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
      providedOpBundlePath := ""
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
        providedOpBundlePath,
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
      providedOpArgs := map[string]string{}
      providedCorrelationId := ""
      providedOpBundlePath := ""
      providedOpName := ""
      providedOpNamespace := ""
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
        providedOpNamespace,
        *providedLogger,
      )

      /* assert */
      receivedCorrelationId,
      receivedArgs,
      receivedOpBundlePath,
      receivedOpName,
      receivedOpNamespace,
      receivedLogger := fakeRunOpUseCase.ExecuteArgsForCall(0)
      Expect(receivedArgs).To(Equal(providedOpArgs))
      Expect(receivedCorrelationId).To(Equal(providedCorrelationId))
      Expect(receivedOpBundlePath).To(Equal(providedOpBundlePath))
      Expect(receivedOpName).To(Equal(providedOpName))
      Expect(receivedOpNamespace).To(Equal(providedOpNamespace))
      Expect(receivedLogger).To(Equal(*providedLogger))
      Expect(fakeRunOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
})
