package osfilesys

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("compositionRoot", func() {
  Context("listNamesOfDevOpDirsUcExecuter", func() {
    It("should return an instance of type listNamesOfDevOpDirsUcExecuter", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualListNamesOfDevOpDirsUcExecuter := objectUnderTest.ListNamesOfDevOpDirsUcExecuter()

      /* assert */
      Expect(actualListNamesOfDevOpDirsUcExecuter).To(BeAssignableToTypeOf(&listNamesOfDevOpDirsUcExecuterImpl{}))

    })
  })
  Context("listNamesOfPipelineDirsUcExecuter", func() {
    It("should return an instance of type listNamesOfPipelineDirsUcExecuter", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualListNamesOfPipelineDirsUcExecuter := objectUnderTest.ListNamesOfPipelineDirsUcExecuter()

      /* assert */
      Expect(actualListNamesOfPipelineDirsUcExecuter).To(BeAssignableToTypeOf(&listNamesOfPipelineDirsUcExecuterImpl{}))

    })
  })
  Context("readDevOpFileUcExecuter", func() {
    It("should return an instance of type readDevOpFileUcExecuter", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualReadDevOpFileUcExecuter := objectUnderTest.ReadDevOpFileUcExecuter()

      /* assert */
      Expect(actualReadDevOpFileUcExecuter).To(BeAssignableToTypeOf(&readDevOpFileUcExecuterImpl{}))

    })
  })
  Context("readPipelineFileUcExecuter", func() {
    It("should return an instance of type readPipelineFileUcExecuter", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualReadPipelineFileUcExecuter := objectUnderTest.ReadPipelineFileUcExecuter()

      /* assert */
      Expect(actualReadPipelineFileUcExecuter).To(BeAssignableToTypeOf(&readPipelineFileUcExecuterImpl{}))

    })
  })
  Context("saveDevOpFileUcExecuter", func() {
    It("should return an instance of type saveDevOpFileUcExecuter", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualSaveDevOpFileUcExecuter := objectUnderTest.SaveDevOpFileUcExecuter()

      /* assert */
      Expect(actualSaveDevOpFileUcExecuter).To(BeAssignableToTypeOf(&saveDevOpFileUcExecuterImpl{}))

    })
  })
  Context("savePipelineFileUcExecuter", func() {
    It("should return an instance of type savePipelineFileUcExecuter", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualSavePipelineFileUcExecuter := objectUnderTest.SavePipelineFileUcExecuter()

      /* assert */
      Expect(actualSavePipelineFileUcExecuter).To(BeAssignableToTypeOf(&savePipelineFileUcExecuterImpl{}))

    })
  })
  Context("createDevOpDirUcExecuter", func() {
    It("should return an instance of type createDevOpDirUcExecuter", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualCreateDevOpDirUcExecuter := objectUnderTest.CreateDevOpDirUcExecuter()

      /* assert */
      Expect(actualCreateDevOpDirUcExecuter).To(BeAssignableToTypeOf(&createDevOpDirUcExecuterImpl{}))

    })
  })
  Context("createPipelineDirUcExecuter", func() {
    It("should return an instance of type createPipelineDirUcExecuter", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualCreatePipelineDirUcExecuter := objectUnderTest.CreatePipelineDirUcExecuter()

      /* assert */
      Expect(actualCreatePipelineDirUcExecuter).To(BeAssignableToTypeOf(&createPipelineDirUcExecuterImpl{}))

    })
  })
})