package opspec

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/adapters"
)

var _ = Describe("_compositionRoot", func() {

  var fakeFilesystem = new(FakeFilesystem)
  var fakeEngineHost = new(adapters.FakeEngineHost)

  Context("CreateCollectionUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeEngineHost, fakeFilesystem)

      /* act */
      actualCreateCollectionUseCase := objectUnderTest.CreateCollectionUseCase()

      /* assert */
      Expect(actualCreateCollectionUseCase).NotTo(BeNil())

    })
  })
  Context("CreateOpUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeEngineHost, fakeFilesystem)

      /* act */
      actualCreateOpUseCase := objectUnderTest.CreateOpUseCase()

      /* assert */
      Expect(actualCreateOpUseCase).NotTo(BeNil())

    })
  })
  Context("GetEventStreamUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeEngineHost, fakeFilesystem)

      /* act */
      actualGetEventStreamUseCase := objectUnderTest.GetEventStreamUseCase()

      /* assert */
      Expect(actualGetEventStreamUseCase).NotTo(BeNil())

    })
  })
  Context("GetCollectionUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeEngineHost, fakeFilesystem)

      /* act */
      actualGetCollectionUseCase := objectUnderTest.GetCollectionUseCase()

      /* assert */
      Expect(actualGetCollectionUseCase).NotTo(BeNil())

    })
  })
  Context("GetOpUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeEngineHost, fakeFilesystem)

      /* act */
      actualGetOpUseCase := objectUnderTest.GetOpUseCase()

      /* assert */
      Expect(actualGetOpUseCase).NotTo(BeNil())

    })
  })
  Context("KillOpRunUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeEngineHost, fakeFilesystem)

      /* act */
      actualKillOpRunUseCase := objectUnderTest.KillOpRunUseCase()

      /* assert */
      Expect(actualKillOpRunUseCase).NotTo(BeNil())

    })
  })
  Context("SetCollectionDescriptionUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeEngineHost, fakeFilesystem)

      /* act */
      actualSetCollectionDescriptionUseCase := objectUnderTest.SetCollectionDescriptionUseCase()

      /* assert */
      Expect(actualSetCollectionDescriptionUseCase).NotTo(BeNil())

    })
  })
  Context("SetOpDescriptionUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeEngineHost, fakeFilesystem)

      /* act */
      actualSetOpDescriptionUseCase := objectUnderTest.SetOpDescriptionUseCase()

      /* assert */
      Expect(actualSetOpDescriptionUseCase).NotTo(BeNil())

    })
  })
  Context("StartOpRunUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeEngineHost, fakeFilesystem)

      /* act */
      actualStartOpRunUseCase := objectUnderTest.StartOpRunUseCase()

      /* assert */
      Expect(actualStartOpRunUseCase).NotTo(BeNil())

    })
  })
  Context("TryResolveDefaultCollectionUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot(fakeEngineHost, fakeFilesystem)

      /* act */
      actualTryResolveDefaultCollectionUseCase := objectUnderTest.TryResolveDefaultCollectionUseCase()

      /* assert */
      Expect(actualTryResolveDefaultCollectionUseCase).NotTo(BeNil())

    })
  })
})
