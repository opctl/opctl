package models

import (
  . "github.com/onsi/ginkgo"
  "fmt"
)

var _ = Describe("OpRunStartedEvent", func() {
  Context("an instance", func() {
    It("should implement models.Event interface", func() {

      /* arrange */
      var objectUnderTest Event

      /* act/assert */
      objectUnderTest = OpRunStartedEvent{}
      fmt.Sprint(objectUnderTest)

    })
  })
})
