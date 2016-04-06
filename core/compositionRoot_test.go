package core

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  fakeContainerEngine "github.com/dev-op-spec/engine/core/adapters/containerengine/fake"
  fakeFilesys "github.com/dev-op-spec/engine/core/adapters/filesys/fake"
)

var _ = Describe("compositionRoot", func() {
  Context("AddOpUseCase()", func() {
    It("should return an instance of type addOpUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualAddOpUseCase := objectUnderTest.AddOpUseCase()

      /* assert */
      Expect(actualAddOpUseCase).To(BeAssignableToTypeOf(&_addOpUseCase{}))

    })
  })
  Context("AddSubOpUseCase()", func() {
    It("should return an instance of type addSubOpUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualAddSubOpUseCase := objectUnderTest.AddSubOpUseCase()

      /* assert */
      Expect(actualAddSubOpUseCase).To(BeAssignableToTypeOf(&_addSubOpUseCase{}))

    })
  })
  Context("ListOpsUseCase()", func() {
    It("should return an instance of type listOpsUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualListOpsUseCase := objectUnderTest.ListOpsUseCase()

      /* assert */
      Expect(actualListOpsUseCase).To(BeAssignableToTypeOf(&_listOpsUseCase{}))

    })
  })
  Context("RunOpUseCase()", func() {
    It("should return an instance of type runOpUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualRunOpUseCase := objectUnderTest.RunOpUseCase()

      /* assert */
      Expect(actualRunOpUseCase).To(BeAssignableToTypeOf(&_runOpUseCase{}))

    })
  })
  Context("SetDescriptionOfOpUseCase()", func() {
    It("should return an instance of type setDescriptionOfOpUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot(
        fakeContainerEngine.New(),
        fakeFilesys.New(),
      )

      /* act */
      actualSetDescriptionOfOpUseCase := objectUnderTest.SetDescriptionOfOpUseCase()

      /* assert */
      Expect(actualSetDescriptionOfOpUseCase).To(BeAssignableToTypeOf(&_setDescriptionOfOpUseCase{}))

    })
  })
})