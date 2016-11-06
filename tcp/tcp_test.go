package tcp

import (
  . "github.com/onsi/ginkgo"
  "github.com/opspec-io/engine/core"
)

var _ = Describe("tcp", func() {
  Context("New", func() {
    It("should return an instance of Tcp", func() {

      /* arrange */
      var _ = New(new(core.FakeCore)).(Api)

    })
  })
  Context("Start", func() {
    It("should not panic", func() {

      /* arrange */
      objectUnderTest := New(new(core.FakeCore))

      /* arrange/act/assert */
      go objectUnderTest.Start()

    })
  })
})
