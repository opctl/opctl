package os

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("compositionRoot", func() {
  Context("CreateDirUseCase", func() {
    It("should return an instance of type createDirUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualCreateDirUseCase := objectUnderTest.CreateDirUseCase()

      /* assert */
      Expect(actualCreateDirUseCase).To(BeAssignableToTypeOf(&_createDirUseCase{}))

    })
  })
  Context("GetBytesOfFileUseCase", func() {
    It("should return an instance of type getBytesOfFileUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualGetBytesOfFileUseCase := objectUnderTest.GetBytesOfFileUseCase()

      /* assert */
      Expect(actualGetBytesOfFileUseCase).To(BeAssignableToTypeOf(&_getBytesOfFileUseCase{}))

    })
  })
  Context("ListNamesOfChildDirsUseCase", func() {
    It("should return an instance of type listNamesOfChildDirsUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualListNamesOfChildDirsUseCase := objectUnderTest.ListNamesOfChildDirsUseCase()

      /* assert */
      Expect(actualListNamesOfChildDirsUseCase).To(BeAssignableToTypeOf(&_listNamesOfChildDirsUseCase{}))

    })
  })
  Context("SaveFileUseCase", func() {
    It("should return an instance of type saveFileUseCase", func() {

      /* arrange */
      objectUnderTest := newCompositionRoot()

      /* act */
      actualSaveFileUseCase := objectUnderTest.SaveFileUseCase()

      /* assert */
      Expect(actualSaveFileUseCase).To(BeAssignableToTypeOf(&_saveFileUseCase{}))

    })
  })
})