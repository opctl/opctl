package opspec

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
  Context("SetCollectionDescriptionUseCase", func() {
    It("should return an instance of type _setCollectionDescriptionUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeFilesystem)

      /* act */
      actualSetCollectionDescriptionUseCase := objectUnderTest.SetCollectionDescriptionUseCase()

      /* assert */
      Expect(actualSetCollectionDescriptionUseCase).To(BeAssignableToTypeOf(&_setCollectionDescriptionUseCase{}))

    })
  })
  Context("SetOpDescriptionUseCase", func() {
    It("should return an instance of type _setOpDescriptionUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeFilesystem)

      /* act */
      actualSetOpDescriptionUseCase := objectUnderTest.SetOpDescriptionUseCase()

      /* assert */
      Expect(actualSetOpDescriptionUseCase).To(BeAssignableToTypeOf(&_setOpDescriptionUseCase{}))

    })
  })
})
