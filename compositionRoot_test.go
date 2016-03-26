package main

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("compositionRoot", func() {
  Context("RestApi()", func() {
    It("should not return nil", func() {

      /* arrange */
      objectUnderTest, _ := newCompositionRoot()

      /* act */
      actualRestApi := objectUnderTest.RestApi()

      /* assert */
      Expect(actualRestApi).ShouldNot(BeNil())

    })
  })
})