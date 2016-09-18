package docker

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("containerEngine", func() {
  Context(".StartContainer() method", func() {
    It("should invoke compositionRoot.startContainerUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedOpArgs := map[string]string{}
      providedOpBundlePath := ""
      providedOpName := ""
      providedOpRunId := ""
      providedRootOpRunId := ""

      fakeStartContainerUseCase := new(fakeStartContainerUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.StartContainerUseCaseReturns(fakeStartContainerUseCase)

      objectUnderTest := &_containerEngine{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.StartContainer(
        providedOpArgs,
        providedOpBundlePath,
        providedOpName,
        providedOpRunId,
        nil,
        providedRootOpRunId,
      )

      /* assert */
      receivedArgs,
      receivedOpBundlePath,
      receivedOpName,
      receivedOpRunId,
      receivedEventPublisher,
      receivedRootOpRunId := fakeStartContainerUseCase.ExecuteArgsForCall(0)
      Expect(receivedArgs).To(Equal(providedOpArgs))
      Expect(receivedOpBundlePath).To(Equal(providedOpBundlePath))
      Expect(receivedOpName).To(Equal(providedOpName))
      Expect(receivedOpRunId).To(Equal(providedOpRunId))
      Expect(receivedEventPublisher).To(BeNil())
      Expect(receivedRootOpRunId).To(Equal(providedRootOpRunId))
      Expect(fakeStartContainerUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
})
