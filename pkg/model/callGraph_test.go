package model

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/util/format"
)

var _ = Describe("CallGraph", func() {
  yaml := format.NewYamlFormat()

  Context("when formatting to/from yaml", func() {

    Context("with non-nil $.op", func() {

      It("should have expected attributes", func() {

        /* arrange */
        expectedCallGraph := CallGraph{
          Op:&OpCall{
            Ref:"dummyOpRef",
            Args:map[string]string{
              "dummyArg1Name":"dummyArg1Value",
            },
            Results:map[string]string{
              "dummyResult1Name":"dummyResult1Value",
            },
          },
        }

        /* act */
        providedYaml, err := yaml.From(expectedCallGraph)
        if (nil != err) {
          panic(err)
        }

        actualCallGraph := CallGraph{}
        yaml.To(providedYaml, &actualCallGraph)

        /* assert */
        Expect(actualCallGraph).To(Equal(expectedCallGraph))

      })

    })

    Context("with non-empty $.parallel", func() {

      It("should have expected attributes", func() {

        /* arrange */
        expectedCallGraph := CallGraph{
          Parallel:[]*CallGraph{
            {
              Op:&OpCall{
                Ref:"dummyOpRef",
              },
            },
          },
        }

        /* act */
        providedYaml, err := yaml.From(expectedCallGraph)
        if (nil != err) {
          panic(err)
        }

        actualCallGraph := CallGraph{}
        yaml.To(providedYaml, &actualCallGraph)

        /* assert */
        Expect(actualCallGraph).To(Equal(expectedCallGraph))

      })

    })
  })

  Context("with non-empty $.serial", func() {

    It("should have expected attributes", func() {

      /* arrange */
      expectedCallGraph := CallGraph{
        Serial:[]*CallGraph{
          {
            Op:&OpCall{
              Ref:"dummyOpRef",
            },
          },
        },
      }

      /* act */
      providedYaml, err := yaml.From(expectedCallGraph)
      if (nil != err) {
        panic(err)
      }

      actualCallGraph := CallGraph{}
      yaml.To(providedYaml, &actualCallGraph)

      /* assert */
      Expect(actualCallGraph).To(Equal(expectedCallGraph))

    })

  })
})
