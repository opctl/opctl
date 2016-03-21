package git

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("compositionRoot", func() {
  Context("getTemplateUcExecuter", func() {
    It("should return an instance of type getTemplateUcExecuter", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualGetTemplateUcExecuter := objectUnderTest.GetTemplateUcExecuter()

      /* assert */
      Expect(actualGetTemplateUcExecuter).To(BeAssignableToTypeOf(&getTemplateUcExecuterImpl{}))

    })
  })
})