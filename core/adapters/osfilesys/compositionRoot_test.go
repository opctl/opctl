package osfilesys

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("compositionRoot", func() {
  Context("listNamesOfDevOpDirsUseCase", func() {
    It("should return an instance of type listNamesOfDevOpDirsUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualListNamesOfDevOpDirsUseCase := objectUnderTest.ListNamesOfDevOpDirsUseCase()

      /* assert */
      Expect(actualListNamesOfDevOpDirsUseCase).To(BeAssignableToTypeOf(&_listNamesOfDevOpDirsUseCase{}))

    })
  })
  Context("listNamesOfPipelineDirsUseCase", func() {
    It("should return an instance of type listNamesOfPipelineDirsUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualListNamesOfPipelineDirsUseCase := objectUnderTest.ListNamesOfPipelineDirsUseCase()

      /* assert */
      Expect(actualListNamesOfPipelineDirsUseCase).To(BeAssignableToTypeOf(&_listNamesOfPipelineDirsUseCase{}))

    })
  })
  Context("readDevOpFileUseCase", func() {
    It("should return an instance of type readDevOpFileUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualReadDevOpFileUseCase := objectUnderTest.ReadDevOpFileUseCase()

      /* assert */
      Expect(actualReadDevOpFileUseCase).To(BeAssignableToTypeOf(&_readDevOpFileUseCase{}))

    })
  })
  Context("readPipelineFileUseCase", func() {
    It("should return an instance of type readPipelineFileUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualReadPipelineFileUseCase := objectUnderTest.ReadPipelineFileUseCase()

      /* assert */
      Expect(actualReadPipelineFileUseCase).To(BeAssignableToTypeOf(&_readPipelineFileUseCase{}))

    })
  })
  Context("saveDevOpFileUseCase", func() {
    It("should return an instance of type saveDevOpFileUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualSaveDevOpFileUseCase := objectUnderTest.SaveDevOpFileUseCase()

      /* assert */
      Expect(actualSaveDevOpFileUseCase).To(BeAssignableToTypeOf(&_saveDevOpFileUseCase{}))

    })
  })
  Context("savePipelineFileUseCase", func() {
    It("should return an instance of type savePipelineFileUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualSavePipelineFileUseCase := objectUnderTest.SavePipelineFileUseCase()

      /* assert */
      Expect(actualSavePipelineFileUseCase).To(BeAssignableToTypeOf(&_savePipelineFileUseCase{}))

    })
  })
  Context("createDevOpDirUseCase", func() {
    It("should return an instance of type createDevOpDirUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualCreateDevOpDirUseCase := objectUnderTest.CreateDevOpDirUseCase()

      /* assert */
      Expect(actualCreateDevOpDirUseCase).To(BeAssignableToTypeOf(&_createDevOpDirUseCase{}))

    })
  })
  Context("createPipelineDirUseCase", func() {
    It("should return an instance of type createPipelineDirUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualCreatePipelineDirUseCase := objectUnderTest.CreatePipelineDirUseCase()

      /* assert */
      Expect(actualCreatePipelineDirUseCase).To(BeAssignableToTypeOf(&_createPipelineDirUseCase{}))

    })
  })
})