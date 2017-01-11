package appdata

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "os/user"
  "fmt"
)

var _ = Describe("appdata", func() {
  Context("GlobalPath", func() {
    It("should return expected path", func() {
      /* arrange */
      expected := "/Library/Application Support"

      objectUnderTest := New()

      /* act */
      result := objectUnderTest.GlobalPath()

      /* assert */
      Expect(result).To(Equal(expected))
    })
  })
  Context("UserPath", func() {
    It("should return expected path", func() {
      /* arrange */
      currentUser, err := user.Current()
      if (nil != err) {
        panic(err)
      }
      expected := fmt.Sprintf("%v/Library/Application Support", currentUser.HomeDir)

      objectUnderTest := New()

      /* act */
      result := objectUnderTest.PerUserPath()

      /* assert */
      Expect(result).To(Equal(expected))
    })
  })
})
