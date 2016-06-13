package dockercompose

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("compositionRoot", func() {

  Context("InitOpUseCase", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot()

      /* act */
      actualInitOpUseCase := objectUnderTest.InitOpUseCase()

      /* assert */
      Expect(actualInitOpUseCase).NotTo(BeNil())

    })
  })

  Context("KillOpRunUseCase", func() {
    It("should return an instance of type killOpRunUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot()

      /* act */
      actualKillOpRunUseCase := objectUnderTest.KillOpRunUseCase()

      /* assert */
      Expect(actualKillOpRunUseCase).NotTo(BeNil())

    })
  })

  Context("RunOpUseCase", func() {
    It("should return an instance of type runOpUseCase", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot()

      /* act */
      actualRunOpUseCase := objectUnderTest.RunOpUseCase()

      /* assert */
      Expect(actualRunOpUseCase).NotTo(BeNil())

    })
  })

})
