package number

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("Validate", func() {
	objectUnderTest := newValidator()
	Context("value nil", func() {
		It("should return expected errors", func() {

			/* arrange */
			expectedErrors := []error{
				errors.New("number required"),
			}

			objectUnderTest := newValidator()

			/* act */
			actualErrors := objectUnderTest.Validate(
				nil,
				&model.NumberConstraints{},
			)

			/* assert */
			Expect(actualErrors).To(Equal(expectedErrors))

		})
	})
	Context("value not nil", func() {
		Context("AllOf constraint", func() {
			Context("value meets all AllOf constraints", func() {

				It("returns no errors", func() {

					/* arrange */
					providedValue := float64(1)
					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
						&model.NumberConstraints{
							AllOf: []*model.NumberConstraints{
								{
									Maximum: providedValue,
								},
							},
						},
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value doesn't meet all AllOf constraints", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValue := float64(2)
					providedConstraints := &model.NumberConstraints{
						AllOf: []*model.NumberConstraints{
							{
								Maximum: providedValue - 1,
							},
						},
					}

					expectedErrors := []error{
						fmt.Errorf(
							`Must be less than or equal to %v`,
							providedConstraints.AllOf[0].Maximum,
						),
						errors.New("Must validate all the schemas (allOf)"),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
						providedConstraints,
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
		Context("AnyOf constraint", func() {
			Context("value meets an AnyOf constraint", func() {

				It("returns no errors", func() {

					/* arrange */
					providedValue := float64(4)
					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
						&model.NumberConstraints{
							AnyOf: []*model.NumberConstraints{
								{
									Maximum: providedValue,
								},
							},
						})

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value doesn't meet an AnyOf constraint", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValue := float64(2)

					providedConstraints := &model.NumberConstraints{
						AnyOf: []*model.NumberConstraints{
							{
								Minimum: providedValue + 1,
							},
						},
					}

					expectedErrors := []error{
						errors.New("Must validate at least one schema (anyOf)"),
						fmt.Errorf(
							`Must be greater than or equal to %v`,
							providedConstraints.AnyOf[0].Minimum,
						),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
						providedConstraints,
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
		Context("Enum constraint", func() {
			Context("value in enum", func() {

				It("returns no errors", func() {

					/* arrange */
					providedValue := float64(4)

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
						&model.NumberConstraints{
							Enum: []float64{providedValue},
						},
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value not in enum", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValue := float64(7.2)
					providedConstraints := &model.NumberConstraints{
						Enum: []float64{providedValue + 1},
					}

					expectedErrors := []error{
						fmt.Errorf(
							`must be one of the following: %v`,
							providedConstraints.Enum[0],
						),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
						providedConstraints,
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
		Context("Format constraint", func() {
			Context("integer", func() {
				Context("value doesn't match Format", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValue := float64(3.3)
						providedConstraints := &model.NumberConstraints{
							Format: "integer",
						}

						expectedErrors := []error{
							fmt.Errorf(
								"Does not match format '%v'",
								providedConstraints.Format,
							),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(
							&providedValue,
							providedConstraints,
						)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value matches Format", func() {

					It("should return no errors", func() {

						/* arrange */
						providedValue := float64(1)
						providedConstraints := &model.NumberConstraints{
							Format: "integer",
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Validate(
							&providedValue,
							providedConstraints,
						)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
		})
		Context("Maximum constraint", func() {
			Context("value == Maximum", func() {

				It("returns no errors", func() {

					/* arrange */
					providedValue := float64(2)
					providedConstraints := &model.NumberConstraints{
						Maximum: providedValue,
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
						providedConstraints,
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value > Maximum", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValue := float64(2)
					providedConstraints := &model.NumberConstraints{
						Maximum: providedValue - 1,
					}

					expectedErrors := []error{
						fmt.Errorf(
							"Must be less than or equal to %v",
							providedConstraints.Maximum,
						),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
						providedConstraints,
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value < Maximum", func() {

				It("returns no errors", func() {

					/* arrange */
					providedValue := float64(1)
					providedConstraints := &model.NumberConstraints{
						Maximum: providedValue + 1,
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
						providedConstraints,
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
		Context("Minimum constraint", func() {
			Context("value == Minimum", func() {

				It("should return no errors", func() {

					/* arrange */
					providedValue := float64(1)
					providedConstraints := &model.NumberConstraints{
						Minimum: providedValue,
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
						providedConstraints,
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value < Minimum", func() {

				It("should return expected errors", func() {

					/* arrange */
					providedValue := float64(1)
					providedConstraints := &model.NumberConstraints{
						Minimum: providedValue + 1,
					}

					expectedErrors := []error{
						fmt.Errorf(
							"Must be greater than or equal to %v",
							providedConstraints.Minimum,
						),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
						providedConstraints,
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value > Minimum", func() {

				It("should return no errors", func() {

					/* arrange */
					providedValue := float64(1)
					providedConstraints := &model.NumberConstraints{
						Minimum: providedValue - 1,
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
						providedConstraints,
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
		Context("Not constraint", func() {
			Context("value matches", func() {

				It("should return expected errors", func() {

					/* arrange */
					providedValue := float64(1)
					providedConstraints := &model.NumberConstraints{
						Not: &model.NumberConstraints{
							Enum: []float64{providedValue},
						},
					}

					expectedErrors := []error{
						errors.New("Must not validate the schema (not)"),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
						providedConstraints,
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value doesn't match", func() {

				It("should return no errors", func() {

					/* arrange */
					providedValue := float64(1)
					providedConstraints := &model.NumberConstraints{
						Not: &model.NumberConstraints{
							Enum: []float64{providedValue - 1},
						},
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
						providedConstraints,
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
		Context("OneOf constraint", func() {
			Context("value meets a single OneOf constraint", func() {

				It("returns no errors", func() {

					/* arrange */
					providedValue := float64(1)
					providedConstraints := &model.NumberConstraints{
						OneOf: []*model.NumberConstraints{
							{
								Maximum: providedValue,
							},
						},
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
						providedConstraints,
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value meets no OneOf constraints", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValue := float64(1)
					providedConstraints := &model.NumberConstraints{
						OneOf: []*model.NumberConstraints{
							{
								Minimum: providedValue + 1,
							},
						},
					}

					expectedErrors := []error{
						errors.New("Must validate one and only one schema (oneOf)"),
						fmt.Errorf(
							`Must be greater than or equal to %v`,
							providedConstraints.OneOf[0].Minimum,
						),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
						providedConstraints,
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value meets multiple OneOf constraints", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValue := float64(4)
					providedConstraints := &model.NumberConstraints{
						OneOf: []*model.NumberConstraints{
							{
								Minimum: providedValue,
							},
							{
								Enum: []float64{providedValue},
							},
						},
					}

					expectedErrors := []error{
						errors.New("Must validate one and only one schema (oneOf)"),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
						providedConstraints,
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
	})
})
