package model

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/util/format"
)

var _ = Describe("Param", func() {

	Context("when formatting to/from yaml", func() {
		yaml := format.NewYamlFormat()

		Context("with non-nil $.dir", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedParam := Param{
					Dir: &DirParam{
						Description: "dummyDescription",
						IsSecret:    true,
					},
				}

				/* act */
				providedYaml, err := yaml.From(expectedParam)
				if nil != err {
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
					File: &FileParam{
						Description: "dummyDescription",
						IsSecret:    true,
					},
				}

				/* act */
				providedYaml, err := yaml.From(expectedParam)
				if nil != err {
					panic(err)
				}

				actualParam := Param{}
				yaml.To(providedYaml, &actualParam)

				/* assert */
				Expect(actualParam).To(Equal(expectedParam))

			})

		})

		Context("with non-nil $.number", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedParam := Param{
					Number: &NumberParam{
						Default:     "dummyDefault",
						Description: "dummyParamDescription",
						Constraints: &NumberConstraints{
							AllOf: []*NumberConstraints{
								{
									Maximum: 2,
								},
							},
							AnyOf: []*NumberConstraints{
								{
									Minimum: 1,
								},
							},
							Enum:       []float64{1.2, 2},
							Integer:    true,
							Maximum:    1000,
							MultipleOf: 1,
							Minimum:    0,
							OneOf: []*NumberConstraints{
								{
									Minimum: 0,
								},
							},
						},
						IsSecret: true,
					},
				}

				/* act */
				providedYaml, err := yaml.From(expectedParam)
				if nil != err {
					panic(err)
				}

				actualParam := Param{}
				yaml.To(providedYaml, &actualParam)

				/* assert */
				Expect(actualParam).To(Equal(expectedParam))

			})

		})

		Context("with non-nil $.socket", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedParam := Param{
					Socket: &SocketParam{
						Description: "dummyDescription",
						IsSecret:    true,
					},
				}

				/* act */
				providedYaml, err := yaml.From(expectedParam)
				if nil != err {
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
					String: &StringParam{
						Default:     "dummyDefault",
						Description: "dummyParamDescription",
						Constraints: &StringConstraints{
							AllOf: []*StringConstraints{
								{
									MaxLength: 2,
								},
							},
							AnyOf: []*StringConstraints{
								{
									Pattern: "xyz",
								},
							},
							Enum:      []string{"dummyEnumItem1", "dummyEnumItem2"},
							Format:    "dummyFormat",
							MaxLength: 1000,
							MinLength: 0,
							OneOf: []*StringConstraints{
								{
									MinLength: 0,
								},
							},
							Pattern: "dummyPattern",
						},
						IsSecret: true,
					},
				}

				/* act */
				providedYaml, err := yaml.From(expectedParam)
				if nil != err {
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
