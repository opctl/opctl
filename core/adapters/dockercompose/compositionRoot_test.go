package dockercompose

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("compositionRoot", func() {
  Context("initDevOpUseCase", func() {
    It("should return an instance of type initDevOpUseCase", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualInitDevOpUseCase := objectUnderTest.InitDevOpUseCase()

      /* assert */
      Expect(actualInitDevOpUseCase).To(BeAssignableToTypeOf(&_initDevOpUseCase{}))

    })
  })
  Context("runDevOpUseCase", func() {
    It("should return an instance of type runDevOpUseCase", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualRunDevOpUseCase := objectUnderTest.RunDevOpUseCase()

      /* assert */
      Expect(actualRunDevOpUseCase).To(BeAssignableToTypeOf(&_runDevOpUseCase{}))

    })
  })
})