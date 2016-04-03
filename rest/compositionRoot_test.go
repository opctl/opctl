package rest

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/dev-op-spec/engine/core"
)

var _ = Describe("compositionRoot", func() {
  Context("AddOperationHandler()", func() {
    It("should return an addOperationHandler instance", func() {

      /* arrange */

      fakeCoreApi := new(core.FakeApi)

      objectUnderTest := newCompositionRoot(fakeCoreApi)

      /* act */
      actualAddOperationHandler := objectUnderTest.AddOperationHandler()

      /* assert */
      Expect(actualAddOperationHandler).To(BeAssignableToTypeOf(&addOperationHandler{}))

    })
  })
})