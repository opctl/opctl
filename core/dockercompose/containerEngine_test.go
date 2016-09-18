package dockercompose

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/engine/core"
)

var _ = Describe("containerEngine", func() {
  Context(".StartOpRun() method", func() {
    It("should invoke compositionRoot.startOpRunUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedOpArgs := map[string]string{}
      providedOpBundlePath := ""
      providedOpName := ""
      providedOpRunId := ""
      providedEventPublisher := new(core.EventPublisher)
      providedRootOpRunId := ""

      fakeStartOpRunUseCase := new(fakeStartOpRunUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.StartOpRunUseCaseReturns(fakeStartOpRunUseCase)

      objectUnderTest := &_containerEngine{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.StartContainer(
        providedOpArgs,
        providedOpBundlePath,
        providedOpName,
        providedOpRunId,
        *providedEventPublisher,
        providedRootOpRunId,
      )

      /* assert */
      receivedArgs,
      receivedOpBundlePath,
      receivedOpName,
      receivedOpRunId,
      receivedEventPublisher,
      receivedRootOpRunId := fakeStartOpRunUseCase.ExecuteArgsForCall(0)
      Expect(receivedArgs).To(Equal(providedOpArgs))
      Expect(receivedOpBundlePath).To(Equal(providedOpBundlePath))
      Expect(receivedOpName).To(Equal(providedOpName))
      Expect(receivedOpRunId).To(Equal(providedOpRunId))
      Expect(receivedEventPublisher).To(Equal(*providedEventPublisher))
      Expect(receivedRootOpRunId).To(Equal(providedRootOpRunId))
      Expect(fakeStartOpRunUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
})
