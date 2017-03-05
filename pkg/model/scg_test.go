package model

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/util/format"
)

var _ = Describe("Scg", func() {

	Context("when formatting to/from yaml", func() {
		yaml := format.NewYamlFormat()

		Context("with non-nil $.op", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedCallGraph := Scg{
					Op: &ScgOpCall{
						Ref: "dummyOpPkgRef",
						Inputs: map[string]string{
							"dummyArg1Name": "dummyArg1Value",
						},
						Outputs: map[string]string{
							"dummyResult1Name": "dummyResult1Value",
						},
					},
				}

				/* act */
				providedYaml, err := yaml.From(expectedCallGraph)
				if nil != err {
					panic(err)
				}

				actualCallGraph := Scg{}
				yaml.To(providedYaml, &actualCallGraph)

				/* assert */
				Expect(actualCallGraph).To(Equal(expectedCallGraph))

			})

		})

		Context("with non-empty $.parallel", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedCallGraph := Scg{
					Parallel: []*Scg{
						{
							Op: &ScgOpCall{
								Ref: "dummyOpPkgRef",
							},
						},
					},
				}

				/* act */
				providedYaml, err := yaml.From(expectedCallGraph)
				if nil != err {
					panic(err)
				}

				actualCallGraph := Scg{}
				yaml.To(providedYaml, &actualCallGraph)

				/* assert */
				Expect(actualCallGraph).To(Equal(expectedCallGraph))

			})

		})

		Context("with non-empty $.serial", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedCallGraph := Scg{
					Serial: []*Scg{
						{
							Op: &ScgOpCall{
								Ref: "dummyOpPkgRef",
							},
						},
					},
				}

				/* act */
				providedYaml, err := yaml.From(expectedCallGraph)
				if nil != err {
					panic(err)
				}

				actualCallGraph := Scg{}
				yaml.To(providedYaml, &actualCallGraph)

				/* assert */
				Expect(actualCallGraph).To(Equal(expectedCallGraph))

			})

		})
	})
})
