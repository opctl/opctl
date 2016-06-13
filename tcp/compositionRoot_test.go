package tcp

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opctl/engine/core"
)

var _ = Describe("compositionRoot", func() {

  Context("GetEventStreamHandler()", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(
        new(core.FakeApi),
      )

      /* act */
      actualGetEventStreamHandler := objectUnderTest.GetEventStreamHandler()

      /* assert */
      Expect(actualGetEventStreamHandler).NotTo(BeNil())

    })
  })

  Context("GetLivenessHandler()", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(
        new(core.FakeApi),
      )

      /* act */
      actualGetLivenessHandler := objectUnderTest.GetLivenessHandler()

      /* assert */
      Expect(actualGetLivenessHandler).NotTo(BeNil())

    })
  })

  Context("KillOpRunHandler()", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(
        new(core.FakeApi),
      )

      /* act */
      actualKillOpRunHandler := objectUnderTest.KillOpRunHandler()

      /* assert */
      Expect(actualKillOpRunHandler).NotTo(BeNil())

    })
  })

  Context("RunOpHandler()", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(
        new(core.FakeApi),
      )

      /* act */
      actualRunOpHandler := objectUnderTest.RunOpHandler()

      /* assert */
      Expect(actualRunOpHandler).NotTo(BeNil())

    })
  })

})
