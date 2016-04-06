package dockercompose

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("compositionRoot", func() {
  Context("initOpUseCase", func() {
    It("should return an instance of type initOpUseCase", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualInitOpUseCase := objectUnderTest.InitOpUseCase()

      /* assert */
      Expect(actualInitOpUseCase).To(BeAssignableToTypeOf(&_initOpUseCase{}))

    })
  })
  Context("runOpUseCase", func() {
    It("should return an instance of type runOpUseCase", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualRunOpUseCase := objectUnderTest.RunOpUseCase()

      /* assert */
      Expect(actualRunOpUseCase).To(BeAssignableToTypeOf(&_runOpUseCase{}))

    })
  })
})