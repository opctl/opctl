package core

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("compositionRoot", func() {
  Context("addDevOpUcExecuter", func() {
    It("should return an instance of type addDevOpUcExecuter", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualAddDevOpUcExecuter := objectUnderTest.AddDevOpUcExecuter()

      /* assert */
      Expect(actualAddDevOpUcExecuter).To(BeAssignableToTypeOf(&addDevOpUcExecuterImpl{}))

    })
  })
  Context("addPipelineUcExecuter", func() {
    It("should return an instance of type addPipelineUcExecuter", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualAddPipelineUcExecuter := objectUnderTest.AddPipelineUcExecuter()

      /* assert */
      Expect(actualAddPipelineUcExecuter).To(BeAssignableToTypeOf(&addPipelineUcExecuterImpl{}))

    })
  })
  Context("addStageToPipelineUcExecuter", func() {
    It("should return an instance of type addStageToPipelineUcExecuter", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualAddStageToPipelineUcExecuter := objectUnderTest.AddStageToPipelineUcExecuter()

      /* assert */
      Expect(actualAddStageToPipelineUcExecuter).To(BeAssignableToTypeOf(&addStageToPipelineUcExecuterImpl{}))

    })
  })
  Context("listDevOpsUcExecuter", func() {
    It("should return an instance of type listDevOpsUcExecuter", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualListDevOpsUcExecuter := objectUnderTest.ListDevOpsUcExecuter()

      /* assert */
      Expect(actualListDevOpsUcExecuter).To(BeAssignableToTypeOf(&listDevOpsUcExecuterImpl{}))

    })
  })
  Context("listPipelinesUcExecuter", func() {
    It("should return an instance of type listPipelinesUcExecuter", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualListPipelinesUcExecuter := objectUnderTest.ListPipelinesUcExecuter()

      /* assert */
      Expect(actualListPipelinesUcExecuter).To(BeAssignableToTypeOf(&listPipelinesUcExecuterImpl{}))

    })
  })
  Context("runDevOpUcExecuter", func() {
    It("should return an instance of type runDevOpUcExecuter", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualRunDevOpUcExecuter := objectUnderTest.RunDevOpUcExecuter()

      /* assert */
      Expect(actualRunDevOpUcExecuter).To(BeAssignableToTypeOf(&runDevOpUcExecuterImpl{}))

    })
  })
  Context("runPipelineUcExecuter", func() {
    It("should return an instance of type runPipelineUcExecuter", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualRunPipelineUcExecuter := objectUnderTest.RunPipelineUcExecuter()

      /* assert */
      Expect(actualRunPipelineUcExecuter).To(BeAssignableToTypeOf(&runPipelineUcExecuterImpl{}))

    })
  })
  Context("setDescriptionOfDevOpUcExecuter", func() {
    It("should return an instance of type setDescriptionOfDevOpUcExecuter", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualSetDescriptionOfDevOpUcExecuter := objectUnderTest.SetDescriptionOfDevOpUcExecuter()

      /* assert */
      Expect(actualSetDescriptionOfDevOpUcExecuter).To(BeAssignableToTypeOf(&setDescriptionOfDevOpUcExecuterImpl{}))

    })
  })
  Context("setDescriptionOfPipelineUcExecuter", func() {
    It("should return an instance of type setDescriptionOfPipelineUcExecuter", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualSetDescriptionOfPipelineUcExecuter := objectUnderTest.SetDescriptionOfPipelineUcExecuter()

      /* assert */
      Expect(actualSetDescriptionOfPipelineUcExecuter).To(BeAssignableToTypeOf(&setDescriptionOfPipelineUcExecuterImpl{}))

    })
  })
})