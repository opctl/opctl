package osfilesys

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("filesys", func() {
  Context(".CreateDevOpDir() method", func() {
    It("should invoke compositionRoot.createDevOpDirUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */
      providedDevOpName := ""

      // wire up fakes
      fakeCreateDevOpDirUCExecuter := new(FakeCreateDevOpDirUcExecuter)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.CreateDevOpDirUcExecuterReturns(fakeCreateDevOpDirUCExecuter)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.CreateDevOpDir(providedDevOpName)

      /* assert */
      Expect(fakeCreateDevOpDirUCExecuter.ExecuteArgsForCall(0)).To(Equal(providedDevOpName))
      Expect(fakeCreateDevOpDirUCExecuter.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".CreatePipelineDir() method", func() {
    It("should invoke compositionRoot.createPipelineDirUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */
      providedPipelineName := ""

      // wire up fakes
      fakeCreatePipelineDirUCExecuter := new(FakeCreatePipelineDirUcExecuter)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.CreatePipelineDirUcExecuterReturns(fakeCreatePipelineDirUCExecuter)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.CreatePipelineDir(providedPipelineName)

      /* assert */
      Expect(fakeCreatePipelineDirUCExecuter.ExecuteArgsForCall(0)).To(Equal(providedPipelineName))
      Expect(fakeCreatePipelineDirUCExecuter.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".ListNamesOfDevOpDirs() method", func() {
    It("should invoke compositionRoot.listNamesOfDevOpDirsUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */

      // wire up fakes
      fakeListNamesOfDevOpDirsUCExecuter := new(FakeListNamesOfDevOpDirsUcExecuter)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.ListNamesOfDevOpDirsUcExecuterReturns(fakeListNamesOfDevOpDirsUCExecuter)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.ListNamesOfDevOpDirs()

      /* assert */
      Expect(fakeListNamesOfDevOpDirsUCExecuter.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".ListNamesOfPipelineDirs() method", func() {
    It("should invoke compositionRoot.listNamesOfPipelineDirsUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */

      // wire up fakes
      fakeListNamesOfPipelineDirsUCExecuter := new(FakeListNamesOfPipelineDirsUcExecuter)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.ListNamesOfPipelineDirsUcExecuterReturns(fakeListNamesOfPipelineDirsUCExecuter)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.ListNamesOfPipelineDirs()

      /* assert */
      Expect(fakeListNamesOfPipelineDirsUCExecuter.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".ReadDevOpFile() method", func() {
    It("should invoke compositionRoot.readDevOpFileUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */
      providedDevOpName := ""

      // wire up fakes
      fakeReadDevOpFileUCExecuter := new(FakeReadDevOpFileUcExecuter)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.ReadDevOpFileUcExecuterReturns(fakeReadDevOpFileUCExecuter)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.ReadDevOpFile(providedDevOpName)

      /* assert */
      Expect(fakeReadDevOpFileUCExecuter.ExecuteArgsForCall(0)).To(Equal(providedDevOpName))
      Expect(fakeReadDevOpFileUCExecuter.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".ReadPipelineFile() method", func() {
    It("should invoke compositionRoot.readPipelineFileUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */
      providedPipelineName := ""

      // wire up fakes
      fakeReadPipelineFileUCExecuter := new(FakeReadPipelineFileUcExecuter)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.ReadPipelineFileUcExecuterReturns(fakeReadPipelineFileUCExecuter)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.ReadPipelineFile(providedPipelineName)

      /* assert */
      Expect(fakeReadPipelineFileUCExecuter.ExecuteArgsForCall(0)).To(Equal(providedPipelineName))
      Expect(fakeReadPipelineFileUCExecuter.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".SaveDevOpFile() method", func() {
    It("should invoke compositionRoot.saveDevOpFileUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */
      providedDevOpName := ""
      providedDevOpData := make([]byte, 0)

      // wire up fakes
      fakeSaveDevOpFileUCExecuter := new(FakeSaveDevOpFileUcExecuter)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.SaveDevOpFileUcExecuterReturns(fakeSaveDevOpFileUCExecuter)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.SaveDevOpFile(providedDevOpName, providedDevOpData)

      /* assert */
      executeArg0, executeArg1 := fakeSaveDevOpFileUCExecuter.ExecuteArgsForCall(0)
      Expect(executeArg0).To(Equal(providedDevOpName))
      Expect(executeArg1).To(Equal(providedDevOpData))
      Expect(fakeSaveDevOpFileUCExecuter.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".SavePipelineFile() method", func() {
    It("should invoke compositionRoot.savePipelineFileUcExecuter.Execute() with expected args & return result", func() {

      /* arrange */
      providedPipelineName := ""
      providedPipelineData := make([]byte, 0)

      // wire up fakes
      fakeSavePipelineFileUCExecuter := new(FakeSavePipelineFileUcExecuter)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.SavePipelineFileUcExecuterReturns(fakeSavePipelineFileUCExecuter)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.SavePipelineFile(providedPipelineName, providedPipelineData)

      /* assert */
      executeArg0, executeArg1 := fakeSavePipelineFileUCExecuter.ExecuteArgsForCall(0)
      Expect(executeArg0).To(Equal(providedPipelineName))
      Expect(executeArg1).To(Equal(providedPipelineData))
      Expect(fakeSavePipelineFileUCExecuter.ExecuteCallCount()).To(Equal(1))

    })
  })
})
