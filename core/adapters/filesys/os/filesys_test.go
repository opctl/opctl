package os

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("filesys", func() {
  Context(".CreateDir() method", func() {
    It("should invoke compositionRoot.createDirUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedPathToDir := ""

      // wire up fakes
      fakeCreateDirUseCase := new(fakeCreateDirUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.CreateDirUseCaseReturns(fakeCreateDirUseCase)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.CreateDir(providedPathToDir)

      /* assert */
      Expect(fakeCreateDirUseCase.ExecuteArgsForCall(0)).To(Equal(providedPathToDir))
      Expect(fakeCreateDirUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".ListNamesOfChildDirs() method", func() {
    It("should invoke compositionRoot.listNamesOfChildDirsUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedPathToPERATIONarentDir := ""

      // wire up fakes
      fakeListNamesOfChildDirsUseCase := new(fakeListNamesOfChildDirsUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.ListNamesOfChildDirsUseCaseReturns(fakeListNamesOfChildDirsUseCase)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.ListNamesOfChildDirs(providedPathToPERATIONarentDir)

      /* assert */
      Expect(fakeListNamesOfChildDirsUseCase.ExecuteArgsForCall(0)).To(Equal(providedPathToPERATIONarentDir))
      Expect(fakeListNamesOfChildDirsUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".GetBytesOfFile() method", func() {
    It("should invoke compositionRoot.getBytesOfFileUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedPathToFile := ""

      // wire up fakes
      fakeGetBytesOfFileUseCase := new(fakeGetBytesOfFileUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.GetBytesOfFileUseCaseReturns(fakeGetBytesOfFileUseCase)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.GetBytesOfFile(providedPathToFile)

      /* assert */
      Expect(fakeGetBytesOfFileUseCase.ExecuteArgsForCall(0)).To(Equal(providedPathToFile))
      Expect(fakeGetBytesOfFileUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".SaveFile() method", func() {
    It("should invoke compositionRoot.saveFileUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedPathToFile := ""
      providedBytesOfFile := make([]byte, 0)

      // wire up fakes
      fakeSaveFileUseCase := new(fakeSaveFileUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.SaveFileUseCaseReturns(fakeSaveFileUseCase)

      objectUnderTest := &filesys{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.SaveFile(providedPathToFile, providedBytesOfFile)

      /* assert */
      executeArg0, executeArg1 := fakeSaveFileUseCase.ExecuteArgsForCall(0)
      Expect(executeArg0).To(Equal(providedPathToFile))
      Expect(executeArg1).To(Equal(providedBytesOfFile))
      Expect(fakeSaveFileUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
})
