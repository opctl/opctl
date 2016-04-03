package core

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  fakeContainerEngine "github.com/dev-op-spec/engine/core/adapters/containerengine/fake"
  fakeFilesys "github.com/dev-op-spec/engine/core/adapters/filesys/fake"
)

var _ = Describe("compositionRoot", func() {
  Context("AddOperationUseCase()", func() {
    It("should return an instance of type addOperationUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualAddOperationUseCase := objectUnderTest.AddOperationUseCase()

      /* assert */
      Expect(actualAddOperationUseCase).To(BeAssignableToTypeOf(&_addOperationUseCase{}))

    })
  })
  Context("AddSubOperationUseCase()", func() {
    It("should return an instance of type addSubOperationUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualAddSubOperationUseCase := objectUnderTest.AddSubOperationUseCase()

      /* assert */
      Expect(actualAddSubOperationUseCase).To(BeAssignableToTypeOf(&_addSubOperationUseCase{}))

    })
  })
  Context("ListOperationsUseCase()", func() {
    It("should return an instance of type listOperationsUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualListOperationsUseCase := objectUnderTest.ListOperationsUseCase()

      /* assert */
      Expect(actualListOperationsUseCase).To(BeAssignableToTypeOf(&_listOperationsUseCase{}))

    })
  })
  Context("RunOperationUseCase()", func() {
    It("should return an instance of type runOperationUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualRunOperationUseCase := objectUnderTest.RunOperationUseCase()

      /* assert */
      Expect(actualRunOperationUseCase).To(BeAssignableToTypeOf(&_runOperationUseCase{}))

    })
  })
  Context("SetDescriptionOfOperationUseCase()", func() {
    It("should return an instance of type setDescriptionOfOperationUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualSetDescriptionOfOperationUseCase := objectUnderTest.SetDescriptionOfOperationUseCase()

      /* assert */
      Expect(actualSetDescriptionOfOperationUseCase).To(BeAssignableToTypeOf(&_setDescriptionOfOperationUseCase{}))

    })
  })
})