package model

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/util/format"
)

var _ = Describe("Arg", func() {

  Context("when formatting to/from json", func() {
    json := format.NewJsonFormat()

    Context("with non-nil $.dir", func() {

      It("should have expected attributes", func() {

        /* arrange */
        expectedArg := Arg{
          Dir:"dummyDirRef",
        }

        /* act */
        providedJson, err := json.From(expectedArg)
        if (nil != err) {
          panic(err)
        }

        actualArg := Arg{}
        json.To(providedJson, &actualArg)

        /* assert */
        Expect(actualArg).To(Equal(expectedArg))

      })

    })

    Context("with non-nil $.file", func() {

      It("should have expected attributes", func() {

        /* arrange */
        expectedArg := Arg{
          File:"dummyFileRef",
        }

        /* act */
        providedJson, err := json.From(expectedArg)
        if (nil != err) {
          panic(err)
        }

        actualArg := Arg{}
        json.To(providedJson, &actualArg)

        /* assert */
        Expect(actualArg).To(Equal(expectedArg))

      })

    })

    Context("with non-nil $.netSocket", func() {

      It("should have expected attributes", func() {

        /* arrange */
        expectedArg := Arg{
          NetSocket:&NetSocketArg{
            Host:"dummyName",
            Port:1,
          },
        }

        /* act */
        providedJson, err := json.From(expectedArg)
        if (nil != err) {
          panic(err)
        }

        actualArg := Arg{}
        json.To(providedJson, &actualArg)

        /* assert */
        Expect(actualArg).To(Equal(expectedArg))

      })

    })

    Context("with non-nil $.string", func() {

      It("should have expected attributes", func() {

        /* arrange */
        expectedArg := Arg{
          String: "dummyString",
        }

        /* act */
        providedJson, err := json.From(expectedArg)
        if (nil != err) {
          panic(err)
        }

        actualArg := Arg{}
        json.To(providedJson, &actualArg)

        /* assert */
        Expect(actualArg).To(Equal(expectedArg))

      })

    })

  })

})
