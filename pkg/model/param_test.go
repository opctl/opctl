package model

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/util/format"
)

var _ = Describe("Param", func() {
  yaml := format.NewYamlFormat()

  Context("when formatting to/from yaml", func() {

    Context("with nil $.dir, $.file, $.netSocket, and $.string", func() {

      It("should have expected attributes", func() {

        /* arrange */
        expectedParam := Param{
          V0_1_2Param: V0_1_2Param{
            Name:"dummyName",
            Description:"dummyDescription",
            IsSecret:true,
            Default:"dummyDefault",
          },
        }

        /* act */
        providedYaml, err := yaml.From(expectedParam)
        if (nil != err) {
          panic(err)
        }

        actualParam := Param{}
        yaml.To(providedYaml, &actualParam)

        /* assert */
        Expect(actualParam).To(Equal(expectedParam))

      })

    })

    Context("with non-nil $.dir", func() {

      It("should have expected attributes", func() {

        /* arrange */
        expectedParam := Param{
          Dir:&DirParam{
            Name:"dummyName",
            Description:"dummyDescription",
            IsSecret:true,
          },
        }

        /* act */
        providedYaml, err := yaml.From(expectedParam)
        if (nil != err) {
          panic(err)
        }

        actualParam := Param{}
        yaml.To(providedYaml, &actualParam)

        /* assert */
        Expect(actualParam).To(Equal(expectedParam))

      })

    })

    Context("with non-nil $.file", func() {

      It("should have expected attributes", func() {

        /* arrange */
        expectedParam := Param{
          File:&FileParam{
            Name:"dummyName",
            Description:"dummyDescription",
            IsSecret:true,
          },
        }

        /* act */
        providedYaml, err := yaml.From(expectedParam)
        if (nil != err) {
          panic(err)
        }

        actualParam := Param{}
        yaml.To(providedYaml, &actualParam)

        /* assert */
        Expect(actualParam).To(Equal(expectedParam))

      })

    })

    Context("with non-nil $.netSocket", func() {

      It("should have expected attributes", func() {

        /* arrange */
        expectedParam := Param{
          NetSocket:&NetSocketParam{
            Name:"dummyName",
            Description:"dummyDescription",
            IsSecret:true,
          },
        }

        /* act */
        providedYaml, err := yaml.From(expectedParam)
        if (nil != err) {
          panic(err)
        }

        actualParam := Param{}
        yaml.To(providedYaml, &actualParam)

        /* assert */
        Expect(actualParam).To(Equal(expectedParam))

      })

    })

    Context("with non-nil $.string", func() {

      It("should have expected attributes", func() {

        /* arrange */
        expectedParam := Param{
          String:&StringParam{
            Name:"dummyName",
            Description:"dummyDescription",
            IsSecret:true,
            Default:"dummyDefault",
          },
        }

        /* act */
        providedYaml, err := yaml.From(expectedParam)
        if (nil != err) {
          panic(err)
        }

        actualParam := Param{}
        yaml.To(providedYaml, &actualParam)

        /* assert */
        Expect(actualParam).To(Equal(expectedParam))

      })

    })

  })

})
