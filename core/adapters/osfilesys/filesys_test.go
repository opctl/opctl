package osfilesys

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("filesys", func() {
  Context(".CreateDevOpDir() method", func() {
    It("should invoke compositionRoot.createDevOpDirUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedDevOpName := ""

      // wire up fakes
      fakeCreateDevOpDirUseCase := new(FakeCreateDevOpDirUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.CreateDevOpDirUseCaseReturns(fakeCreateDevOpDirUseCase)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.CreateDevOpDir(providedDevOpName)

      /* assert */
      Expect(fakeCreateDevOpDirUseCase.ExecuteArgsForCall(0)).To(Equal(providedDevOpName))
      Expect(fakeCreateDevOpDirUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".CreatePipelineDir() method", func() {
    It("should invoke compositionRoot.createPipelineDirUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedPipelineName := ""

      // wire up fakes
      fakeCreatePipelineDirUseCase := new(FakeCreatePipelineDirUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.CreatePipelineDirUseCaseReturns(fakeCreatePipelineDirUseCase)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.CreatePipelineDir(providedPipelineName)

      /* assert */
      Expect(fakeCreatePipelineDirUseCase.ExecuteArgsForCall(0)).To(Equal(providedPipelineName))
      Expect(fakeCreatePipelineDirUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".ListNamesOfDevOpDirs() method", func() {
    It("should invoke compositionRoot.listNamesOfDevOpDirsUseCase.Execute() with expected args & return result", func() {

      /* arrange */

      // wire up fakes
      fakeListNamesOfDevOpDirsUseCase := new(FakeListNamesOfDevOpDirsUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.ListNamesOfDevOpDirsUseCaseReturns(fakeListNamesOfDevOpDirsUseCase)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.ListNamesOfDevOpDirs()

      /* assert */
      Expect(fakeListNamesOfDevOpDirsUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".ListNamesOfPipelineDirs() method", func() {
    It("should invoke compositionRoot.listNamesOfPipelineDirsUseCase.Execute() with expected args & return result", func() {

      /* arrange */

      // wire up fakes
      fakeListNamesOfPipelineDirsUseCase := new(FakeListNamesOfPipelineDirsUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.ListNamesOfPipelineDirsUseCaseReturns(fakeListNamesOfPipelineDirsUseCase)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.ListNamesOfPipelineDirs()

      /* assert */
      Expect(fakeListNamesOfPipelineDirsUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".ReadDevOpFile() method", func() {
    It("should invoke compositionRoot.readDevOpFileUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedDevOpName := ""

      // wire up fakes
      fakeReadDevOpFileUseCase := new(FakeReadDevOpFileUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.ReadDevOpFileUseCaseReturns(fakeReadDevOpFileUseCase)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.ReadDevOpFile(providedDevOpName)

      /* assert */
      Expect(fakeReadDevOpFileUseCase.ExecuteArgsForCall(0)).To(Equal(providedDevOpName))
      Expect(fakeReadDevOpFileUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".ReadPipelineFile() method", func() {
    It("should invoke compositionRoot.readPipelineFileUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedPipelineName := ""

      // wire up fakes
      fakeReadPipelineFileUseCase := new(FakeReadPipelineFileUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.ReadPipelineFileUseCaseReturns(fakeReadPipelineFileUseCase)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.ReadPipelineFile(providedPipelineName)

      /* assert */
      Expect(fakeReadPipelineFileUseCase.ExecuteArgsForCall(0)).To(Equal(providedPipelineName))
      Expect(fakeReadPipelineFileUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".SaveDevOpFile() method", func() {
    It("should invoke compositionRoot.saveDevOpFileUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedDevOpName := ""
      providedDevOpData := make([]byte, 0)

      // wire up fakes
      fakeSaveDevOpFileUseCase := new(FakeSaveDevOpFileUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.SaveDevOpFileUseCaseReturns(fakeSaveDevOpFileUseCase)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.SaveDevOpFile(providedDevOpName, providedDevOpData)

      /* assert */
      executeArg0, executeArg1 := fakeSaveDevOpFileUseCase.ExecuteArgsForCall(0)
      Expect(executeArg0).To(Equal(providedDevOpName))
      Expect(executeArg1).To(Equal(providedDevOpData))
      Expect(fakeSaveDevOpFileUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".SavePipelineFile() method", func() {
    It("should invoke compositionRoot.savePipelineFileUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedPipelineName := ""
      providedPipelineData := make([]byte, 0)

      // wire up fakes
      fakeSavePipelineFileUseCase := new(FakeSavePipelineFileUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.SavePipelineFileUseCaseReturns(fakeSavePipelineFileUseCase)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.SavePipelineFile(providedPipelineName, providedPipelineData)

      /* assert */
      executeArg0, executeArg1 := fakeSavePipelineFileUseCase.ExecuteArgsForCall(0)
      Expect(executeArg0).To(Equal(providedPipelineName))
      Expect(executeArg1).To(Equal(providedPipelineData))
      Expect(fakeSavePipelineFileUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
})
