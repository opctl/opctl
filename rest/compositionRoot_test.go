package rest

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/dev-op-spec/engine/core"
)

var _ = Describe("compositionRoot", func() {
  Context("CoreApi", func() {
    It("should return a core.Api instance", func() {

      /* arrange */

      fakeCoreApi := new(core.FakeApi)

      objectUnderTest := newCompositionRoot(fakeCoreApi)

      /* act */
      actualCoreApi := objectUnderTest.CoreApi()

      /* assert */
      Expect(actualCoreApi).To(Equal(fakeCoreApi))

    })
  })
})