package paramvalidator

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Describe("Validate", func() {
	objectUnderTest := New()
	Context("invoked w/ non-nil param.Number", func() {
		Context("& non-zero value.Number", func() {
			Context("AllOf constraint", func() {
				Context("value meets all AllOf constraints", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValueNumber := float64(1)
						providedValue := &model.Data{
							Number: &providedValueNumber,
						}
						providedParam := &model.Param{
							Number: &model.NumberParam{
								Constraints: &model.NumberConstraints{
									AllOf: []*model.NumberConstraints{
										{
											Maximum: *providedValue.Number,
										},
									},
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value doesn't meet all AllOf constraints", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValueNumber := float64(2)
						providedValue := &model.Data{
							Number: &providedValueNumber,
						}
						providedParam := &model.Param{
							Number: &model.NumberParam{
								Constraints: &model.NumberConstraints{
									AllOf: []*model.NumberConstraints{
										{
											Maximum: *providedValue.Number - 1,
										},
									},
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf(
								`Must be less than or equal to %v`,
								providedParam.Number.Constraints.AllOf[0].Maximum,
							),
							errors.New("Must validate all the schemas (allOf)"),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
			Context("AnyOf constraint", func() {
				Context("value meets an AnyOf constraint", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValueNumber := float64(4)
						providedValue := &model.Data{
							Number: &providedValueNumber,
						}
						providedParam := &model.Param{
							Number: &model.NumberParam{
								Constraints: &model.NumberConstraints{
									AnyOf: []*model.NumberConstraints{
										{
											Maximum: *providedValue.Number,
										},
									},
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value doesn't meet an AnyOf constraint", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValueNumber := float64(2)
						providedValue := &model.Data{
							Number: &providedValueNumber,
						}
						providedParam := &model.Param{
							Number: &model.NumberParam{
								Constraints: &model.NumberConstraints{
									AnyOf: []*model.NumberConstraints{
										{
											Minimum: *providedValue.Number + 1,
										},
									},
								},
							},
						}

						expectedErrors := []error{
							errors.New("Must validate at least one schema (anyOf)"),
							fmt.Errorf(
								`Must be greater than or equal to %v`,
								providedParam.Number.Constraints.AnyOf[0].Minimum,
							),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
			Context("Enum constraint", func() {
				Context("value in enum", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValueNumber := float64(4)
						providedValue := &model.Data{
							Number: &providedValueNumber,
						}
						providedParam := &model.Param{
							Number: &model.NumberParam{
								Constraints: &model.NumberConstraints{
									Enum: []float64{*providedValue.Number},
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value not in enum", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValueNumber := float64(7.2)
						providedValue := &model.Data{
							Number: &providedValueNumber,
						}
						providedParam := &model.Param{
							Number: &model.NumberParam{
								Constraints: &model.NumberConstraints{
									Enum: []float64{*providedValue.Number + 1},
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf(
								`must be one of the following: %v`,
								providedParam.Number.Constraints.Enum[0],
							),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

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
							providedValueNumber := float64(3.3)
							providedValue := &model.Data{
								Number: &providedValueNumber,
							}
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										Format: "integer",
									},
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"Does not match format '%v'",
									providedParam.Number.Constraints.Format,
								),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("value matches Format", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValueNumber := float64(1)
							providedValue := &model.Data{
								Number: &providedValueNumber,
							}
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										Format: "integer",
									},
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

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
						providedValueNumber := float64(2)
						providedValue := &model.Data{
							Number: &providedValueNumber,
						}
						providedParam := &model.Param{
							Number: &model.NumberParam{
								Constraints: &model.NumberConstraints{
									Maximum: *providedValue.Number,
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value > Maximum", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValueNumber := float64(2)
						providedValue := &model.Data{
							Number: &providedValueNumber,
						}
						providedParam := &model.Param{
							Number: &model.NumberParam{
								Constraints: &model.NumberConstraints{
									Maximum: *providedValue.Number - 1,
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf(
								"Must be less than or equal to %v",
								providedParam.Number.Constraints.Maximum,
							),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value < Maximum", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValueNumber := float64(1)
						providedValue := &model.Data{
							Number: &providedValueNumber,
						}
						providedParam := &model.Param{
							Number: &model.NumberParam{
								Constraints: &model.NumberConstraints{
									Maximum: *providedValue.Number + 1,
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
			Context("Minimum constraint", func() {
				Context("value == Minimum", func() {

					It("should return no errors", func() {

						/* arrange */
						providedValueNumber := float64(1)
						providedValue := &model.Data{
							Number: &providedValueNumber,
						}
						providedParam := &model.Param{
							Number: &model.NumberParam{
								Constraints: &model.NumberConstraints{
									Minimum: *providedValue.Number,
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value < Minimum", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValueNumber := float64(1)
						providedValue := &model.Data{
							Number: &providedValueNumber,
						}
						providedParam := &model.Param{
							Number: &model.NumberParam{
								Constraints: &model.NumberConstraints{
									Minimum: *providedValue.Number + 1,
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf(
								"Must be greater than or equal to %v",
								providedParam.Number.Constraints.Minimum,
							),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value > Minimum", func() {

					It("should return no errors", func() {

						/* arrange */
						providedValueNumber := float64(1)
						providedValue := &model.Data{
							Number: &providedValueNumber,
						}
						providedParam := &model.Param{
							Number: &model.NumberParam{
								Constraints: &model.NumberConstraints{
									Minimum: *providedValue.Number - 1,
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
			Context("Not constraint", func() {
				Context("value matches", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValueNumber := float64(1)
						providedValue := &model.Data{
							Number: &providedValueNumber,
						}
						providedParam := &model.Param{
							Number: &model.NumberParam{
								Constraints: &model.NumberConstraints{
									Not: &model.NumberConstraints{
										Enum: []float64{*providedValue.Number},
									},
								},
							},
						}

						expectedErrors := []error{
							errors.New("Must not validate the schema (not)"),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value doesn't match", func() {

					It("should return no errors", func() {

						/* arrange */
						providedValueNumber := float64(1)
						providedValue := &model.Data{
							Number: &providedValueNumber,
						}
						providedParam := &model.Param{
							Number: &model.NumberParam{
								Constraints: &model.NumberConstraints{
									Not: &model.NumberConstraints{
										Enum: []float64{*providedValue.Number - 1},
									},
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
			Context("OneOf constraint", func() {
				Context("value meets a single OneOf constraint", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValueNumber := float64(1)
						providedValue := &model.Data{
							Number: &providedValueNumber,
						}
						providedParam := &model.Param{
							Number: &model.NumberParam{
								Constraints: &model.NumberConstraints{
									OneOf: []*model.NumberConstraints{
										{
											Maximum: *providedValue.Number,
										},
									},
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value meets no OneOf constraints", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValueNumber := float64(1)
						providedValue := &model.Data{
							Number: &providedValueNumber,
						}
						providedParam := &model.Param{
							Number: &model.NumberParam{
								Constraints: &model.NumberConstraints{
									OneOf: []*model.NumberConstraints{
										{
											Minimum: *providedValue.Number + 1,
										},
									},
								},
							},
						}

						expectedErrors := []error{
							errors.New("Must validate one and only one schema (oneOf)"),
							fmt.Errorf(
								`Must be greater than or equal to %v`,
								providedParam.Number.Constraints.OneOf[0].Minimum,
							),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value meets multiple OneOf constraints", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValueNumber := float64(4)
						providedValue := &model.Data{
							Number: &providedValueNumber,
						}
						providedParam := &model.Param{
							Number: &model.NumberParam{
								Constraints: &model.NumberConstraints{
									OneOf: []*model.NumberConstraints{
										{
											Minimum: *providedValue.Number,
										},
										{
											Enum: []float64{*providedValue.Number},
										},
									},
								},
							},
						}

						expectedErrors := []error{
							errors.New("Must validate one and only one schema (oneOf)"),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
		})
		Context("& nil value.Number", func() {
			Context("and non empty Default", func() {
				Context("AllOf constraint", func() {
					Context("default meets all AllOf constraints", func() {

						It("returns no errors", func() {

							/* arrange */
							providedValue := &model.Data{}
							providedDefault := float64(3)
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										AllOf: []*model.NumberConstraints{
											{
												Maximum: providedDefault,
											},
										},
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("value doesn't meet all AllOf constraints", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedValue := &model.Data{}
							providedDefault := float64(3)
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										AllOf: []*model.NumberConstraints{
											{
												Minimum: providedDefault + 1,
											},
										},
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									`Must be greater than or equal to %v`,
									providedParam.Number.Constraints.AllOf[0].Minimum,
								),
								errors.New("Must validate all the schemas (allOf)"),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("AnyOf constraint", func() {
					Context("default meets an AnyOf constraint", func() {

						It("returns no errors", func() {

							/* arrange */
							providedValue := &model.Data{}
							providedDefault := float64(3)
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										AnyOf: []*model.NumberConstraints{
											{
												Maximum: providedDefault,
											},
										},
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("value doesn't meet an AnyOf constraint", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedValue := &model.Data{}
							providedDefault := float64(3)
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										AnyOf: []*model.NumberConstraints{
											{
												Enum: []float64{providedDefault + 1},
											},
										},
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{
								errors.New("Must validate at least one schema (anyOf)"),
								fmt.Errorf(
									`must be one of the following: %v`,
									providedParam.Number.Constraints.AnyOf[0].Enum[0],
								),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("Enum constraint", func() {
					Context("default in enum", func() {

						It("returns no errors", func() {

							/* arrange */
							providedValue := &model.Data{}
							providedDefault := float64(3)
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										Enum: []float64{providedDefault},
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default not in enum", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedValue := &model.Data{}
							providedDefault := float64(3)
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										Enum: []float64{providedDefault + 1},
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									`must be one of the following: %v`,
									providedParam.Number.Constraints.Enum[0],
								),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("Maximum constraint", func() {
					Context("default == Maximum", func() {

						It("returns no errors", func() {

							/* arrange */
							providedValue := &model.Data{}
							providedDefault := float64(3)
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										Maximum: providedDefault,
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default < Maximum", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedValue := &model.Data{}
							providedDefault := float64(3)
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										Maximum: providedDefault - 1,
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"Must be less than or equal to %v",
									providedParam.Number.Constraints.Maximum,
								),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default > Maximum", func() {

						It("returns no errors", func() {

							/* arrange */
							providedValue := &model.Data{}
							providedDefault := float64(3)
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										Maximum: providedDefault + 1,
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("Minimum constraint", func() {
					Context("default == Minimum", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValue := &model.Data{}
							providedDefault := float64(3)
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										Minimum: providedDefault,
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default < Minimum", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValue := &model.Data{}
							providedDefault := float64(3)
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										Minimum: providedDefault + 1,
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"Must be greater than or equal to %v",
									providedParam.Number.Constraints.Minimum,
								),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default > Minimum", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValue := &model.Data{}
							providedDefault := float64(3)
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										Minimum: providedDefault - 1,
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("Not constraint", func() {
					Context("default matches", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValue := &model.Data{}
							providedDefault := float64(3)
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										Not: &model.NumberConstraints{
											Minimum: providedDefault,
										},
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{
								errors.New("Must not validate the schema (not)"),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default doesn't match", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValue := &model.Data{}
							providedDefault := float64(3)
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										Not: &model.NumberConstraints{
											Maximum: providedDefault - 1,
										},
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("OneOf constraint", func() {
					Context("default meets a single OneOf constraint", func() {

						It("returns no errors", func() {

							/* arrange */
							providedValue := &model.Data{}
							providedDefault := float64(3)
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										OneOf: []*model.NumberConstraints{
											{
												Minimum: providedDefault,
											},
										},
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default meets no OneOf constraints", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedValue := &model.Data{}
							providedDefault := float64(3)
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										OneOf: []*model.NumberConstraints{
											{
												Enum: []float64{providedDefault + 1},
											},
										},
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{
								errors.New("Must validate one and only one schema (oneOf)"),
								fmt.Errorf(
									`must be one of the following: %v`,
									providedParam.Number.Constraints.OneOf[0].Enum[0],
								),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default meets multiple OneOf constraints", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedValue := &model.Data{}
							providedDefault := float64(3)
							providedParam := &model.Param{
								Number: &model.NumberParam{
									Constraints: &model.NumberConstraints{
										OneOf: []*model.NumberConstraints{
											{
												Minimum: providedDefault,
											},
											{
												Enum: []float64{providedDefault},
											},
										},
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{
								errors.New("Must validate one and only one schema (oneOf)"),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
			})
		})
	})
})
