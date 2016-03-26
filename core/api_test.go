package core

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
"github.com/dev-op-spec/engine/core/models"
)

var _ = Describe("_sdk", func() {
  Context(".AddDevOp() method", func() {
    It("should invoke compositionRoot.addDevOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedAddDevOpReq := models.NewAddDevOpReq("", "")

      // wire up fakes
      fakeAddDevOpUseCase := new(fakeAddDevOpUseCase)
      
      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.AddDevOpUseCaseReturns(fakeAddDevOpUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.AddDevOp(*providedAddDevOpReq)

      /* assert */
      Expect(fakeAddDevOpUseCase.ExecuteArgsForCall(0)).To(Equal(*providedAddDevOpReq))
      Expect(fakeAddDevOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".AddPipeline() method", func() {
    It("should invoke compositionRoot.addPipelineUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedAddPipelineReq := models.NewAddPipelineReq("", "")

      // wire up fakes
      fakeAddPipelineUseCase := new(fakeAddPipelineUseCase)
      
      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.AddPipelineUseCaseReturns(fakeAddPipelineUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.AddPipeline(*providedAddPipelineReq)

      /* assert */
      Expect(fakeAddPipelineUseCase.ExecuteArgsForCall(0)).To(Equal(*providedAddPipelineReq))
      Expect(fakeAddPipelineUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".AddStageToPipeline() method", func() {
    It("should invoke compositionRoot.addStageToPipelineUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedAddStageToPipelineReq := models.NewAddStageToPipelineReq(false, "", "", "")

      // wire up fakes
      fakeAddStageToPipelineUseCase := new(fakeAddStageToPipelineUseCase)
      
      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.AddStageToPipelineUseCaseReturns(fakeAddStageToPipelineUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.AddStageToPipeline(*providedAddStageToPipelineReq)

      /* assert */
      Expect(fakeAddStageToPipelineUseCase.ExecuteArgsForCall(0)).To(Equal(*providedAddStageToPipelineReq))
      Expect(fakeAddStageToPipelineUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".ListDevOps() method", func() {
    It("should invoke compositionRoot.listDevOpsUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      expectedReturnedDevOps := make([]models.DevOpView, 0)

      // wire up fakes
      fakeListDevOpsUseCase := new(fakeListDevOpsUseCase)
      fakeListDevOpsUseCase.ExecuteReturns(expectedReturnedDevOps, nil)
      
      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.ListDevOpsUseCaseReturns(fakeListDevOpsUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      actualReturnedDevOps, _ := objectUnderTest.ListDevOps()

      /* assert */
      Expect(actualReturnedDevOps).To(Equal(expectedReturnedDevOps))

    })
  })
  Context(".ListPipelines() method", func() {
    It("should invoke compositionRoot.listPipelinesUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      expectedReturnedPipelines := make([]models.PipelineView, 0)

      // wire up fakes
      fakeListPipelinesUseCase := new(fakeListPipelinesUseCase)
      fakeListPipelinesUseCase.ExecuteReturns(expectedReturnedPipelines, nil)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.ListPipelinesUseCaseReturns(fakeListPipelinesUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      actualReturnedPipelines, _ := objectUnderTest.ListPipelines()

      /* assert */
      Expect(actualReturnedPipelines).To(Equal(expectedReturnedPipelines))

    })
  })
  Context(".RunDevOp() method", func() {
    It("should invoke compositionRoot.runDevOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedDevOpName := ""

      // wire up fakes
      fakeRunDevOpUseCase := new(fakeRunDevOpUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.RunDevOpUseCaseReturns(fakeRunDevOpUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.RunDevOp(providedDevOpName)

      /* assert */
      Expect(fakeRunDevOpUseCase.ExecuteArgsForCall(0)).To(Equal(providedDevOpName))
      Expect(fakeRunDevOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".RunPipeline() method", func() {
    It("should invoke compositionRoot.runPipelineUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedPipelineName := ""

      // wire up fakes
      fakeRunPipelineUseCase := new(fakeRunPipelineUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.RunPipelineUseCaseReturns(fakeRunPipelineUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.RunPipeline(providedPipelineName)

      /* assert */
      executeArg0, executeArg1 := fakeRunPipelineUseCase.ExecuteArgsForCall(0)
      Expect(executeArg0).To(Equal(providedPipelineName))
      Expect(executeArg1).To(Equal(make([]string, 0)))
      Expect(fakeRunPipelineUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".SetDescriptionOfDevOp() method", func() {
    It("should invoke compositionRoot.setDescriptionOfDevOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedSetDescriptionOfDevOpReq := models.NewSetDescriptionOfDevOpReq("", "")

      // wire up fakes
      fakeSetDescriptionOfDevOpUseCase := new(fakeSetDescriptionOfDevOpUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.SetDescriptionOfDevOpUseCaseReturns(fakeSetDescriptionOfDevOpUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.SetDescriptionOfDevOp(*providedSetDescriptionOfDevOpReq)

      /* assert */
      Expect(fakeSetDescriptionOfDevOpUseCase.ExecuteArgsForCall(0)).To(Equal(*providedSetDescriptionOfDevOpReq))
      Expect(fakeSetDescriptionOfDevOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".SetDescriptionOfPipeline() method", func() {
    It("should invoke compositionRoot.setDescriptionOfPipelineUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedSetDescriptionOfPipelineReq := models.NewSetDescriptionOfPipelineReq("", "")

      // wire up fakes
      fakeSetDescriptionOfPipelineUseCase := new(fakeSetDescriptionOfPipelineUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.SetDescriptionOfPipelineUseCaseReturns(fakeSetDescriptionOfPipelineUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.SetDescriptionOfPipeline(*providedSetDescriptionOfPipelineReq)

      /* assert */
      Expect(fakeSetDescriptionOfPipelineUseCase.ExecuteArgsForCall(0)).To(Equal(*providedSetDescriptionOfPipelineReq))
      Expect(fakeSetDescriptionOfPipelineUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })

})
