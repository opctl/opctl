package model

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/util/format"
)

var _ = Describe("CallGraphDeclaration", func() {
  yaml := format.NewYamlFormat()

  Context("when formatting to/from yaml", func() {

    Context("with non-nil $.op", func() {

      It("should have expected attributes", func() {

        /* arrange */
        expectedCallGraphDeclaration := CallGraphDeclaration{
          Op:&OpCallDeclaration{
            Ref:"dummyOpRef",
          },
        }

        /* act */
        providedYaml, err := yaml.From(expectedCallGraphDeclaration)
        if (nil != err) {
          panic(err)
        }

        actualCallGraphDeclaration := CallGraphDeclaration{}
        yaml.To(providedYaml, &actualCallGraphDeclaration)

        /* assert */
        Expect(actualCallGraphDeclaration).To(Equal(expectedCallGraphDeclaration))

      })

    })

    Context("with non-nil $.parallel", func() {

      It("should have expected attributes", func() {

        /* arrange */
        expectedCallGraphDeclaration := CallGraphDeclaration{
          Parallel:&ParallelCallDeclaration{
            {
              Op:&OpCallDeclaration{
                Ref:"dummyOpRef",
              },
            },
          },
        }

        /* act */
        providedYaml, err := yaml.From(expectedCallGraphDeclaration)
        if (nil != err) {
          panic(err)
        }

        actualCallGraphDeclaration := CallGraphDeclaration{}
        yaml.To(providedYaml, &actualCallGraphDeclaration)

        /* assert */
        Expect(actualCallGraphDeclaration).To(Equal(expectedCallGraphDeclaration))

      })

    })
  })

  Context("with non-nil $.serial", func() {

    It("should have expected attributes", func() {

      /* arrange */
      expectedCallGraphDeclaration := CallGraphDeclaration{
        Serial:&SerialCallDeclaration{
          {
            Op:&OpCallDeclaration{
              Ref:"dummyOpRef",
            },
          },
        },
      }

      /* act */
      providedYaml, err := yaml.From(expectedCallGraphDeclaration)
      if (nil != err) {
        panic(err)
      }

      actualCallGraphDeclaration := CallGraphDeclaration{}
      yaml.To(providedYaml, &actualCallGraphDeclaration)

      /* assert */
      Expect(actualCallGraphDeclaration).To(Equal(expectedCallGraphDeclaration))

    })

  })
})
