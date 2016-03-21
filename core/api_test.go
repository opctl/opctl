package core

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
"github.com/dev-op-spec/engine/core/models"
)

var _ = Describe("_sdk", func() {
  Context(".AddDevOp() method", func() {
    It("should invoke compositionRoot.addDevOpUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */
      providedAddDevOpReq := models.NewAddDevOpReq("", "")

      // wire up fakes
      fakeAddDevOpUcExecuter := new(fakeAddDevOpUcExecuter)
      
      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.AddDevOpUcExecuterReturns(fakeAddDevOpUcExecuter)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.AddDevOp(*providedAddDevOpReq)

      /* assert */
      Expect(fakeAddDevOpUcExecuter.ExecuteArgsForCall(0)).To(Equal(*providedAddDevOpReq))
      Expect(fakeAddDevOpUcExecuter.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".AddPipeline() method", func() {
    It("should invoke compositionRoot.addPipelineUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */
      providedAddPipelineReq := models.NewAddPipelineReq("", "")

      // wire up fakes
      fakeAddPipelineUcExecuter := new(fakeAddPipelineUcExecuter)
      
      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.AddPipelineUcExecuterReturns(fakeAddPipelineUcExecuter)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.AddPipeline(*providedAddPipelineReq)

      /* assert */
      Expect(fakeAddPipelineUcExecuter.ExecuteArgsForCall(0)).To(Equal(*providedAddPipelineReq))
      Expect(fakeAddPipelineUcExecuter.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".AddStageToPipeline() method", func() {
    It("should invoke compositionRoot.addStageToPipelineUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */
      providedAddStageToPipelineReq := models.NewAddStageToPipelineReq(false, "", "", "")

      // wire up fakes
      fakeAddStageToPipelineUcExecuter := new(fakeAddStageToPipelineUcExecuter)
      
      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.AddStageToPipelineUcExecuterReturns(fakeAddStageToPipelineUcExecuter)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.AddStageToPipeline(*providedAddStageToPipelineReq)

      /* assert */
      Expect(fakeAddStageToPipelineUcExecuter.ExecuteArgsForCall(0)).To(Equal(*providedAddStageToPipelineReq))
      Expect(fakeAddStageToPipelineUcExecuter.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".ListDevOps() method", func() {
    It("should invoke compositionRoot.listDevOpsUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */
      expectedReturnedDevOps := make([]models.DevOpView, 0)

      // wire up fakes
      fakeListDevOpsUcExecuter := new(fakeListDevOpsUcExecuter)
      fakeListDevOpsUcExecuter.ExecuteReturns(expectedReturnedDevOps, nil)
      
      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.ListDevOpsUcExecuterReturns(fakeListDevOpsUcExecuter)

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
    It("should invoke compositionRoot.listPipelinesUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */
      expectedReturnedPipelines := make([]models.PipelineView, 0)

      // wire up fakes
      fakeListPipelinesUcExecuter := new(fakeListPipelinesUcExecuter)
      fakeListPipelinesUcExecuter.ExecuteReturns(expectedReturnedPipelines, nil)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.ListPipelinesUcExecuterReturns(fakeListPipelinesUcExecuter)

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
    It("should invoke compositionRoot.runDevOpUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */
      providedDevOpName := ""

      // wire up fakes
      fakeRunDevOpUcExecuter := new(fakeRunDevOpUcExecuter)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.RunDevOpUcExecuterReturns(fakeRunDevOpUcExecuter)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.RunDevOp(providedDevOpName)

      /* assert */
      Expect(fakeRunDevOpUcExecuter.ExecuteArgsForCall(0)).To(Equal(providedDevOpName))
      Expect(fakeRunDevOpUcExecuter.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".RunPipeline() method", func() {
    It("should invoke compositionRoot.runPipelineUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */
      providedPipelineName := ""

      // wire up fakes
      fakeRunPipelineUcExecuter := new(fakeRunPipelineUcExecuter)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.RunPipelineUcExecuterReturns(fakeRunPipelineUcExecuter)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.RunPipeline(providedPipelineName)

      /* assert */
      executeArg0, executeArg1 := fakeRunPipelineUcExecuter.ExecuteArgsForCall(0)
      Expect(executeArg0).To(Equal(providedPipelineName))
      Expect(executeArg1).To(Equal(make([]string, 0)))
      Expect(fakeRunPipelineUcExecuter.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".SetDescriptionOfDevOp() method", func() {
    It("should invoke compositionRoot.setDescriptionOfDevOpUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */
      providedSetDescriptionOfDevOpReq := models.NewSetDescriptionOfDevOpReq("", "")

      // wire up fakes
      fakeSetDescriptionOfDevOpUcExecuter := new(fakeSetDescriptionOfDevOpUcExecuter)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.SetDescriptionOfDevOpUcExecuterReturns(fakeSetDescriptionOfDevOpUcExecuter)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.SetDescriptionOfDevOp(*providedSetDescriptionOfDevOpReq)

      /* assert */
      Expect(fakeSetDescriptionOfDevOpUcExecuter.ExecuteArgsForCall(0)).To(Equal(*providedSetDescriptionOfDevOpReq))
      Expect(fakeSetDescriptionOfDevOpUcExecuter.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".SetDescriptionOfPipeline() method", func() {
    It("should invoke compositionRoot.setDescriptionOfPipelineUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */
      providedSetDescriptionOfPipelineReq := models.NewSetDescriptionOfPipelineReq("", "")

      // wire up fakes
      fakeSetDescriptionOfPipelineUcExecuter := new(fakeSetDescriptionOfPipelineUcExecuter)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.SetDescriptionOfPipelineUcExecuterReturns(fakeSetDescriptionOfPipelineUcExecuter)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.SetDescriptionOfPipeline(*providedSetDescriptionOfPipelineReq)

      /* assert */
      Expect(fakeSetDescriptionOfPipelineUcExecuter.ExecuteArgsForCall(0)).To(Equal(*providedSetDescriptionOfPipelineReq))
      Expect(fakeSetDescriptionOfPipelineUcExecuter.ExecuteCallCount()).To(Equal(1))

    })
  })

})
