package core

import (
  . "github.com/onsi/ginkgo"
)

var _ = Describe("uniqueStringFactory", func() {
  Context("Construct()", func() {
    It("should not return the same string in 100000 iterations", func() {

      /* arrange */
      objectUnderTest := newUniqueStringFactory()
      stringsReturnedFromConstruct := map[string]bool{}

      /* act/assert */
      for i := 0; i < 100000; i++ {

        uniqueString, err := objectUnderTest.Construct()
        if (nil != err) {
          Fail(err.Error())
        }

        if _, ok := stringsReturnedFromConstruct[*uniqueString]; ok {
          Fail("same string returned twice")
        }else {
          stringsReturnedFromConstruct[*uniqueString] = true
        }

      }

    })
  })
})
