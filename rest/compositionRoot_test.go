package rest

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/dev-op-spec/engine/core"
)

var _ = Describe("compositionRoot", func() {
  Context("AddDevOpHandler()", func() {
    It("should return an addDevOpHandler instance", func() {

      /* arrange */

      fakeCoreApi := new(core.FakeApi)

      objectUnderTest := newCompositionRoot(fakeCoreApi)

      /* act */
      actualAddDevOpHandler := objectUnderTest.AddDevOpHandler()

      /* assert */
      Expect(actualAddDevOpHandler).To(BeAssignableToTypeOf(&addDevOpHandler{}))

    })
  })
})