package opspec

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("_compositionRoot", func() {

  var fakeFilesystem = new(FakeFilesystem)
  
  Context("CreateCollectionUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeFilesystem)

      /* act */
      actualCreateCollectionUseCase := objectUnderTest.CreateCollectionUseCase()

      /* assert */
      Expect(actualCreateCollectionUseCase).NotTo(BeNil())

    })
  })
  Context("CreateOpUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeFilesystem)

      /* act */
      actualCreateOpUseCase := objectUnderTest.CreateOpUseCase()

      /* assert */
      Expect(actualCreateOpUseCase).NotTo(BeNil())

    })
  })
  Context("GetCollectionUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeFilesystem)

      /* act */
      actualGetCollectionUseCase := objectUnderTest.GetCollectionUseCase()

      /* assert */
      Expect(actualGetCollectionUseCase).NotTo(BeNil())

    })
  })
  Context("GetOpUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeFilesystem)

      /* act */
      actualGetOpUseCase := objectUnderTest.GetOpUseCase()

      /* assert */
      Expect(actualGetOpUseCase).NotTo(BeNil())

    })
  })
  Context("SetCollectionDescriptionUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeFilesystem)

      /* act */
      actualSetCollectionDescriptionUseCase := objectUnderTest.SetCollectionDescriptionUseCase()

      /* assert */
      Expect(actualSetCollectionDescriptionUseCase).NotTo(BeNil())

    })
  })
  Context("SetOpDescriptionUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeFilesystem)

      /* act */
      actualSetOpDescriptionUseCase := objectUnderTest.SetOpDescriptionUseCase()

      /* assert */
      Expect(actualSetOpDescriptionUseCase).NotTo(BeNil())

    })
  })

  Context("TryResolveDefaultCollectionUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeFilesystem)

      /* act */
      actualTryResolveDefaultCollectionUseCase := objectUnderTest.TryResolveDefaultCollectionUseCase()

      /* assert */
      Expect(actualTryResolveDefaultCollectionUseCase).NotTo(BeNil())

    })
  })
})
