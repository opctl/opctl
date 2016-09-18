package docker

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("compositionRoot", func() {

  Context("EnsureContainerRemovedUseCase", func() {
    It("should not be nil", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot()

      /* act */
      actualEnsureContainerRemovedUseCase := objectUnderTest.EnsureContainerRemovedUseCase()

      /* assert */
      Expect(actualEnsureContainerRemovedUseCase).NotTo(BeNil())

    })
  })

  Context("StartContainerUseCase", func() {
    It("should not be nil", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot()

      /* act */
      actualStartContainerUseCase := objectUnderTest.StartContainerUseCase()

      /* assert */
      Expect(actualStartContainerUseCase).NotTo(BeNil())

    })
  })
})
