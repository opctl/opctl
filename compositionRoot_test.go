package main

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("compositionRoot", func() {
  Context("TcpApi()", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest,_ := newCompositionRoot()

      /* act */
      actualTcpApi := objectUnderTest.TcpApi()

      /* assert */
      Expect(actualTcpApi).ShouldNot(BeNil())

    })
  })
})
