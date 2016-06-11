package sdk

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("_compositionRoot", func() {

  var fakeFilesystem = new(FakeFilesystem)

  Context("CreateOpUseCase", func() {
    It("should return an instance of type _createOpUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeFilesystem)

      /* act */
      actualCreateOpUseCase := objectUnderTest.CreateOpUseCase()

      /* assert */
      Expect(actualCreateOpUseCase).To(BeAssignableToTypeOf(&_createOpUseCase{}))

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
