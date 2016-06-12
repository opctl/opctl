package opspec

import (
  . "github.com/onsi/ginkgo"
)

var _ = Describe("_compositionRoot", func() {

  var fakeFilesystem = new(FakeFilesystem)

  Context("CreateOpUseCase", func() {
    It("should return an instance of type createOpUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeFilesystem)

      /* act */
      actualCreateOpUseCase := objectUnderTest.CreateOpUseCase()

      /* assert */
      _, ok := actualCreateOpUseCase.(createOpUseCase)
      if !ok {
        Fail("result not assignable to createOpUseCase")
      }

    })
  })
  Context("SetCollectionDescriptionUseCase", func() {
    It("should return an instance of type setCollectionDescriptionUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeFilesystem)

      /* act */
      actualSetCollectionDescriptionUseCase := objectUnderTest.SetCollectionDescriptionUseCase()

      /* assert */
      _, ok := actualSetCollectionDescriptionUseCase.(setCollectionDescriptionUseCase)
      if !ok {
        Fail("result not assignable to setCollectionDescriptionUseCase")
      }

    })
  })
  Context("SetOpDescriptionUseCase", func() {
    It("should return an instance of type setOpDescriptionUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeFilesystem)

      /* act */
      actualSetOpDescriptionUseCase := objectUnderTest.SetOpDescriptionUseCase()

      /* assert */
      _, ok := actualSetOpDescriptionUseCase.(setOpDescriptionUseCase)
      if !ok {
        Fail("result not assignable to setOpDescriptionUseCase")
      }

    })
  })
})
