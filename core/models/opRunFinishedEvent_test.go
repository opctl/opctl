package models

import (
  . "github.com/onsi/ginkgo"
  "fmt"
)

var _ = Describe("OpRunFinishedEvent", func() {
  Context("an instance", func() {
    It("should implement models.Event interface", func() {

      /* arrange */
      var objectUnderTest Event

      /* act/assert */
      objectUnderTest = opRunFinishedEvent{}
      fmt.Sprint(objectUnderTest)

    })
  })
})
