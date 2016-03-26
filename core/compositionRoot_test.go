package core

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  fakeContainerEngine "github.com/dev-op-spec/engine/core/adapters/containerengine/fake"
  fakeFilesys "github.com/dev-op-spec/engine/core/adapters/filesys/fake"
)

var _ = Describe("compositionRoot", func() {
  Context("AddDevOpUseCase()", func() {
    It("should return an instance of type addDevOpUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualAddDevOpUseCase := objectUnderTest.AddDevOpUseCase()

      /* assert */
      Expect(actualAddDevOpUseCase).To(BeAssignableToTypeOf(&_addDevOpUseCase{}))

    })
  })
  Context("AddPipelineUseCase()", func() {
    It("should return an instance of type addPipelineUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualAddPipelineUseCase := objectUnderTest.AddPipelineUseCase()

      /* assert */
      Expect(actualAddPipelineUseCase).To(BeAssignableToTypeOf(&_addPipelineUseCase{}))

    })
  })
  Context("AddStageToPipelineUseCase()", func() {
    It("should return an instance of type addStageToPipelineUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualAddStageToPipelineUseCase := objectUnderTest.AddStageToPipelineUseCase()

      /* assert */
      Expect(actualAddStageToPipelineUseCase).To(BeAssignableToTypeOf(&_addStageToPipelineUseCase{}))

    })
  })
  Context("ListDevOpsUseCase()", func() {
    It("should return an instance of type listDevOpsUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualListDevOpsUseCase := objectUnderTest.ListDevOpsUseCase()

      /* assert */
      Expect(actualListDevOpsUseCase).To(BeAssignableToTypeOf(&_listDevOpsUseCase{}))

    })
  })
  Context("ListPipelinesUseCase()", func() {
    It("should return an instance of type listPipelinesUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualListPipelinesUseCase := objectUnderTest.ListPipelinesUseCase()

      /* assert */
      Expect(actualListPipelinesUseCase).To(BeAssignableToTypeOf(&_listPipelinesUseCase{}))

    })
  })
  Context("RunDevOpUseCase()", func() {
    It("should return an instance of type runDevOpUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualRunDevOpUseCase := objectUnderTest.RunDevOpUseCase()

      /* assert */
      Expect(actualRunDevOpUseCase).To(BeAssignableToTypeOf(&_runDevOpUseCase{}))

    })
  })
  Context("RunPipelineUseCase()", func() {
    It("should return an instance of type runPipelineUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualRunPipelineUseCase := objectUnderTest.RunPipelineUseCase()

      /* assert */
      Expect(actualRunPipelineUseCase).To(BeAssignableToTypeOf(&_runPipelineUseCase{}))

    })
  })
  Context("SetDescriptionOfDevOpUseCase()", func() {
    It("should return an instance of type setDescriptionOfDevOpUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualSetDescriptionOfDevOpUseCase := objectUnderTest.SetDescriptionOfDevOpUseCase()

      /* assert */
      Expect(actualSetDescriptionOfDevOpUseCase).To(BeAssignableToTypeOf(&_setDescriptionOfDevOpUseCase{}))

    })
  })
  Context("SetDescriptionOfPipelineUseCase()", func() {
    It("should return an instance of type setDescriptionOfPipelineUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualSetDescriptionOfPipelineUseCase := objectUnderTest.SetDescriptionOfPipelineUseCase()

      /* assert */
      Expect(actualSetDescriptionOfPipelineUseCase).To(BeAssignableToTypeOf(&_setDescriptionOfPipelineUseCase{}))

    })
  })
})