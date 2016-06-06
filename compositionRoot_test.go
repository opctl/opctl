package sdk

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("_compositionRoot", func() {

  var fakeFilesystem = new(FakeFilesystem)

  Context("AddOpUseCase", func() {
    It("should return an instance of type _addOpUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeFilesystem)

      /* act */
      actualAddOpUseCase := objectUnderTest.AddOpUseCase()

      /* assert */
      Expect(actualAddOpUseCase).To(BeAssignableToTypeOf(&_addOpUseCase{}))

    })
  })
  Context("SetDescriptionOfOpUseCase", func() {
    It("should return an instance of type _setDescriptionOfOpUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeFilesystem)

      /* act */
      actualSetDescriptionOfOpUseCase := objectUnderTest.SetDescriptionOfOpUseCase()

      /* assert */
      Expect(actualSetDescriptionOfOpUseCase).To(BeAssignableToTypeOf(&_setDescriptionOfOpUseCase{}))

    })
  })
})
