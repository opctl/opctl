package docker

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("compositionRoot", func() {
  Context("EnsureRunningUseCase()", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualEnsureEngineRunningUseCase := objectUnderTest.EnsureEngineRunningUseCase()

      /* assert */
      Expect(actualEnsureEngineRunningUseCase).ShouldNot(BeNil())

    })
  })
  Context("GetEngineProtocolRelativeBaseUrlUseCase()", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualGetEngineProtocolRelativeBaseUrlUseCase := objectUnderTest.GetEngineProtocolRelativeBaseUrlUseCase()

      /* assert */
      Expect(actualGetEngineProtocolRelativeBaseUrlUseCase).ShouldNot(BeNil())

    })
  })
})
