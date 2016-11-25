package model

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/util/format"
)

var _ = Describe("RunDeclaration", func() {
  yaml := format.NewYamlFormat()

  Context("when formatting to/from yaml", func() {

    Context("with non-nil $.op", func() {

      It("should have expected attributes", func() {

        /* arrange */
        expectedRunDeclaration := RunDeclaration{
          Op:"dummyOpRunDeclaration",
        }

        /* act */
        providedYaml, err := yaml.From(expectedRunDeclaration)
        if (nil != err) {
          panic(err)
        }

        actualRunDeclaration := RunDeclaration{}
        yaml.To(providedYaml, &actualRunDeclaration)

        /* assert */
        Expect(actualRunDeclaration).To(Equal(expectedRunDeclaration))

      })

    })

    Context("with non-nil $.parallel", func() {

      It("should have expected attributes", func() {

        /* arrange */
        expectedRunDeclaration := RunDeclaration{
          Parallel:&ParallelRunDeclaration{
            {
              Op:"dummyOpRunDeclaration",
            },
          },
        }

        /* act */
        providedYaml, err := yaml.From(expectedRunDeclaration)
        if (nil != err) {
          panic(err)
        }

        actualRunDeclaration := RunDeclaration{}
        yaml.To(providedYaml, &actualRunDeclaration)

        /* assert */
        Expect(actualRunDeclaration).To(Equal(expectedRunDeclaration))

      })

    })
  })

  Context("with non-nil $.serial", func() {

    It("should have expected attributes", func() {

      /* arrange */
      expectedRunDeclaration := RunDeclaration{
        Serial:&SerialRunDeclaration{
          {
            Op:"dummyOpRunDeclaration",
          },
        },
      }

      /* act */
      providedYaml, err := yaml.From(expectedRunDeclaration)
      if (nil != err) {
        panic(err)
      }

      actualRunDeclaration := RunDeclaration{}
      yaml.To(providedYaml, &actualRunDeclaration)

      /* assert */
      Expect(actualRunDeclaration).To(Equal(expectedRunDeclaration))

    })

  })
})
