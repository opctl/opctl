package dockercompose

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("compositionRoot", func() {
  Context("initDevOpUcExecuter", func() {
    It("should return an instance of type initDevOpUcExecuter", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualInitDevOpUcExecuter := objectUnderTest.InitDevOpUcExecuter()

      /* assert */
      Expect(actualInitDevOpUcExecuter).To(BeAssignableToTypeOf(&initDevOpUcExecuterImpl{}))

    })
  })
  Context("runDevOpUcExecuter", func() {
    It("should return an instance of type runDevOpUcExecuter", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualRunDevOpUcExecuter := objectUnderTest.RunDevOpUcExecuter()

      /* assert */
      Expect(actualRunDevOpUcExecuter).To(BeAssignableToTypeOf(&runDevOpUcExecuterImpl{}))

    })
  })
})