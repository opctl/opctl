package git

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("compositionRoot", func() {
  Context("getTemplateUseCase", func() {
    It("should return an instance of type getTemplateUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualGetTemplateUseCase := objectUnderTest.GetTemplateUseCase()

      /* assert */
      Expect(actualGetTemplateUseCase).To(BeAssignableToTypeOf(&getTemplateUseCaseImpl{}))

    })
  })
})