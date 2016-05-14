package models

import (
  . "github.com/onsi/ginkgo"
  "fmt"
)

var _ = Describe("OpRunKilledEvent", func() {
  Context("an instance", func() {
    It("should implement models.Event interface", func() {

      /* arrange */
      var objectUnderTest Event

      /* act/assert */
      objectUnderTest = opRunKilledEvent{}
      fmt.Sprint(objectUnderTest)

    })
  })
})
