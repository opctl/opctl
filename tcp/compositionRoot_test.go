package tcp

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/dev-op-spec/engine/core"
)

var _ = Describe("compositionRoot", func() {
  Context("AddOpHandler()", func() {
    It("should return an addOpHandler instance", func() {

      /* arrange */

      fakeCoreApi := new(core.FakeApi)

      objectUnderTest := newCompositionRoot(fakeCoreApi)

      /* act */
      actualAddOpHandler := objectUnderTest.AddOpHandler()

      /* assert */
      Expect(actualAddOpHandler).To(BeAssignableToTypeOf(&addOpHandler{}))

    })
  })
})