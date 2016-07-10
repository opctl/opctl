package tcp

import (
  . "github.com/onsi/ginkgo"
  "github.com/opctl/engine/core"
)

var _ = Describe("api", func() {
  Context("New", func() {
    It("should return an instance of Api", func() {

      /* arrange */
      var _ = New(new(core.FakeApi)).(Api)

    })
  })
  Context("Start", func() {
    It("should not panic", func() {

      /* arrange */
      objectUnderTest := New(new(core.FakeApi))

      /* arrange/act/assert */
      go objectUnderTest.Start()

    })
  })
})
