package core

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("compositionRoot", func() {
  Context("addDevOpUseCase", func() {
    It("should return an instance of type addDevOpUseCase", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualAddDevOpUseCase := objectUnderTest.AddDevOpUseCase()

      /* assert */
      Expect(actualAddDevOpUseCase).To(BeAssignableToTypeOf(&_addDevOpUseCase{}))

    })
  })
  Context("addPipelineUseCase", func() {
    It("should return an instance of type addPipelineUseCase", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualAddPipelineUseCase := objectUnderTest.AddPipelineUseCase()

      /* assert */
      Expect(actualAddPipelineUseCase).To(BeAssignableToTypeOf(&_addPipelineUseCase{}))

    })
  })
  Context("addStageToPipelineUseCase", func() {
    It("should return an instance of type addStageToPipelineUseCase", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualAddStageToPipelineUseCase := objectUnderTest.AddStageToPipelineUseCase()

      /* assert */
      Expect(actualAddStageToPipelineUseCase).To(BeAssignableToTypeOf(&_addStageToPipelineUseCase{}))

    })
  })
  Context("listDevOpsUseCase", func() {
    It("should return an instance of type listDevOpsUseCase", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualListDevOpsUseCase := objectUnderTest.ListDevOpsUseCase()

      /* assert */
      Expect(actualListDevOpsUseCase).To(BeAssignableToTypeOf(&_listDevOpsUseCase{}))

    })
  })
  Context("listPipelinesUseCase", func() {
    It("should return an instance of type listPipelinesUseCase", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualListPipelinesUseCase := objectUnderTest.ListPipelinesUseCase()

      /* assert */
      Expect(actualListPipelinesUseCase).To(BeAssignableToTypeOf(&_listPipelinesUseCase{}))

    })
  })
  Context("runDevOpUseCase", func() {
    It("should return an instance of type runDevOpUseCase", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualRunDevOpUseCase := objectUnderTest.RunDevOpUseCase()

      /* assert */
      Expect(actualRunDevOpUseCase).To(BeAssignableToTypeOf(&_runDevOpUseCase{}))

    })
  })
  Context("runPipelineUseCase", func() {
    It("should return an instance of type runPipelineUseCase", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualRunPipelineUseCase := objectUnderTest.RunPipelineUseCase()

      /* assert */
      Expect(actualRunPipelineUseCase).To(BeAssignableToTypeOf(&_runPipelineUseCase{}))

    })
  })
  Context("setDescriptionOfDevOpUseCase", func() {
    It("should return an instance of type setDescriptionOfDevOpUseCase", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualSetDescriptionOfDevOpUseCase := objectUnderTest.SetDescriptionOfDevOpUseCase()

      /* assert */
      Expect(actualSetDescriptionOfDevOpUseCase).To(BeAssignableToTypeOf(&_setDescriptionOfDevOpUseCase{}))

    })
  })
  Context("setDescriptionOfPipelineUseCase", func() {
    It("should return an instance of type setDescriptionOfPipelineUseCase", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualSetDescriptionOfPipelineUseCase := objectUnderTest.SetDescriptionOfPipelineUseCase()

      /* assert */
      Expect(actualSetDescriptionOfPipelineUseCase).To(BeAssignableToTypeOf(&_setDescriptionOfPipelineUseCase{}))

    })
  })
})