package dockercompose

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("compositionRoot", func() {
  Context("initOperationUseCase", func() {
    It("should return an instance of type initOperationUseCase", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualInitOperationUseCase := objectUnderTest.InitOperationUseCase()

      /* assert */
      Expect(actualInitOperationUseCase).To(BeAssignableToTypeOf(&_initOperationUseCase{}))

    })
  })
  Context("runOperationUseCase", func() {
    It("should return an instance of type runOperationUseCase", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualRunOperationUseCase := objectUnderTest.RunOperationUseCase()

      /* assert */
      Expect(actualRunOperationUseCase).To(BeAssignableToTypeOf(&_runOperationUseCase{}))

    })
  })
})