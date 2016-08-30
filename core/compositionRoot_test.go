package core

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/engine/core/adapters/containerengine/fake"
)

var _ = Describe("compositionRoot", func() {

  Context("GetEventStreamUseCase()", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(
        new(fake.FakeContainerEngine),
      )

      /* act */
      actualGetEventStreamUseCase := objectUnderTest.GetEventStreamUseCase()

      /* assert */
      Expect(actualGetEventStreamUseCase).NotTo(BeNil())

    })
  })

  Context("KillOpRunUseCase()", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(
        new(fake.FakeContainerEngine),
      )

      /* act */
      actualKillOpRunUseCase := objectUnderTest.KillOpRunUseCase()

      /* assert */
      Expect(actualKillOpRunUseCase).NotTo(BeNil())

    })
  })

  Context("RunOpUseCase()", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(
        new(fake.FakeContainerEngine),
      )

      /* act */
      actualRunOpUseCase := objectUnderTest.RunOpUseCase()

      /* assert */
      Expect(actualRunOpUseCase).NotTo(BeNil())

    })
  })

})
