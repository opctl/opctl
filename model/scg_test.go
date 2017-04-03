package model

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/util/format"
)

var _ = Describe("SCG", func() {

	Context("when formatting to/from yaml", func() {
		yaml := format.NewYamlFormat()

		Context("with non-nil $.op", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedCallGraph := SCG{
					Op: &SCGOpCall{
						Pkg: &SCGOpCallPkg{
							Ref: "dummyPkgRef",
						},
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

				actualCallGraph := SCG{}
				yaml.To(providedYaml, &actualCallGraph)

				/* assert */
				Expect(actualCallGraph).To(Equal(expectedCallGraph))

			})

		})

		Context("with non-empty $.parallel", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedCallGraph := SCG{
					Parallel: []*SCG{
						{
							Op: &SCGOpCall{
								Pkg: &SCGOpCallPkg{
									Ref: "dummyPkgRef",
								},
							},
						},
					},
				}

				/* act */
				providedYaml, err := yaml.From(expectedCallGraph)
				if nil != err {
					panic(err)
				}

				actualCallGraph := SCG{}
				yaml.To(providedYaml, &actualCallGraph)

				/* assert */
				Expect(actualCallGraph).To(Equal(expectedCallGraph))

			})

		})

		Context("with non-empty $.serial", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedCallGraph := SCG{
					Serial: []*SCG{
						{
							Op: &SCGOpCall{
								Pkg: &SCGOpCallPkg{
									Ref: "dummyPkgRef",
								},
							},
						},
					},
				}

				/* act */
				providedYaml, err := yaml.From(expectedCallGraph)
				if nil != err {
					panic(err)
				}

				actualCallGraph := SCG{}
				yaml.To(providedYaml, &actualCallGraph)

				/* assert */
				Expect(actualCallGraph).To(Equal(expectedCallGraph))

			})

		})
	})
})
