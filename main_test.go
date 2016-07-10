package main

import (
  . "github.com/onsi/ginkgo"
)

var _ = Describe("main", func() {
  Context("main()", func() {
    It("should not panic", func() {

      /* arrange/act/assert */
      go main()

    })
  })
})
