package inputs

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Describe("Validate", func() {
	objectUnderTest := newValidator()
	Context("invoked w/ non-nil param.String", func() {
		Context("& non-empty value.String", func() {
			Context("AllOf constraint", func() {
				Context("value meets all AllOf constraints", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									AllOf: []*model.StringConstraints{
										{
											Pattern: "^.*$",
										},
										{
											Pattern: *providedValue.String,
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
						providedValueString := "dummyValue==\""
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									AllOf: []*model.StringConstraints{
										{
											Pattern: "^$",
										},
										{
											Pattern: *providedValue.String,
										},
									},
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf(
								`Does not match pattern '%v'`,
								providedParam.String.Constraints.AllOf[0].Pattern,
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
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									AnyOf: []*model.StringConstraints{
										{
											Pattern: "^.*$",
										},
										{
											Pattern: *providedValue.String,
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
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									AnyOf: []*model.StringConstraints{
										{
											Pattern: "^$",
										},
										{
											Enum: []string{"dummyEnumItem"},
										},
									},
								},
							},
						}

						expectedErrors := []error{
							errors.New("Must validate at least one schema (anyOf)"),
							fmt.Errorf(
								`Does not match pattern '%v'`,
								providedParam.String.Constraints.AnyOf[0].Pattern,
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
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Enum: []string{*providedValue.String},
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
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Enum: []string{"dummyEnumItem"},
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf(
								`must be one of the following: "%v"`,
								providedParam.String.Constraints.Enum[0],
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
				Context("date-time", func() {
					Context("value doesn't match Format", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValueString := "notDateTime"
							providedValue := &model.Value{
								String: &providedValueString,
							}
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Format: "date-time",
									},
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"Does not match format '%v'",
									providedParam.String.Constraints.Format,
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
							providedValueString := "0000-01-01T00:00:01.0Z"
							providedValue := &model.Value{
								String: &providedValueString,
							}
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Format: "date-time",
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
				Context("docker-image-ref", func() {
					Context("value doesn't match Format", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValueString := "$notADockerImageRef"
							providedValue := &model.Value{
								String: &providedValueString,
							}
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Format: "docker-image-ref",
									},
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"Does not match format '%v'",
									providedParam.String.Constraints.Format,
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
							providedValueString := "dummy-registry.com/dummy-namespace/dummy-repo:dummy-tag"
							providedValue := &model.Value{
								String: &providedValueString,
							}
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Format: "docker-image-ref",
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
				Context("email", func() {
					Context("value doesn't match Format", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValueString := "notEmail"
							providedValue := &model.Value{
								String: &providedValueString,
							}
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Format: "email",
									},
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"Does not match format '%v'",
									providedParam.String.Constraints.Format,
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
							providedValueString := "dummy-email@dummy-domain.com"
							providedValue := &model.Value{
								String: &providedValueString,
							}
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Format: "email",
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
				Context("hostname", func() {
					Context("value doesn't match Format", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValueString := "$notAHostname$"
							providedValue := &model.Value{
								String: &providedValueString,
							}
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Format: "hostname",
									},
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"Does not match format '%v'",
									providedParam.String.Constraints.Format,
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
							providedValueString := "dummy.com"
							providedValue := &model.Value{
								String: &providedValueString,
							}
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Format: "hostname",
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
				Context("ipv4", func() {
					Context("value doesn't match Format", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValueString := "notAnIpV4"
							providedValue := &model.Value{
								String: &providedValueString,
							}
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Format: "ipv4",
									},
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"Does not match format '%v'",
									providedParam.String.Constraints.Format,
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
							providedValueString := "0.0.0.0"
							providedValue := &model.Value{
								String: &providedValueString,
							}
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Format: "ipv4",
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
				Context("ipv6", func() {
					Context("value doesn't match Format", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValueString := "notAnIpV6"
							providedValue := &model.Value{
								String: &providedValueString,
							}
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Format: "ipv6",
									},
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"Does not match format '%v'",
									providedParam.String.Constraints.Format,
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
							providedValueString := "0000:0000:0000:0000:0000:0000:0000:0000"
							providedValue := &model.Value{
								String: &providedValueString,
							}
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Format: "ipv6",
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
				Context("semver", func() {
					Context("value doesn't match Format", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValueString := "$notASemver$"
							providedValue := &model.Value{
								String: &providedValueString,
							}
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Format: "semver",
									},
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"Does not match format '%v'",
									providedParam.String.Constraints.Format,
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
							providedValueString := "1.1.1"
							providedValue := &model.Value{
								String: &providedValueString,
							}
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Format: "semver",
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
				Context("uri", func() {
					Context("value doesn't match Format", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValueString := "notUri"
							providedValue := &model.Value{
								String: &providedValueString,
							}
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Format: "uri",
									},
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"Does not match format '%v'",
									providedParam.String.Constraints.Format,
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
							providedValueString := "https://dummyuri.com:8080/somepath?somequery#somefragment"
							providedValue := &model.Value{
								String: &providedValueString,
							}
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Format: "uri",
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
			Context("MaxLength constraint", func() {
				Context("value length == MaxLength", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									MaxLength: len(*providedValue.String),
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
				Context("value length > MaxLength", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									MaxLength: len(*providedValue.String) - 1,
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf(
								"String length must be less than or equal to %v",
								providedParam.String.Constraints.MaxLength,
							),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value length < MaxLength", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									MaxLength: len(*providedValue.String) + 1,
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
			Context("MinLength constraint", func() {
				Context("value length == MinLength", func() {

					It("should return no errors", func() {

						/* arrange */
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									MinLength: len(*providedValue.String),
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
				Context("value length < MinLength", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									MinLength: len(*providedValue.String) + 1,
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf(
								"String length must be greater than or equal to %v",
								providedParam.String.Constraints.MinLength,
							),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value length > MinLength", func() {

					It("should return no errors", func() {

						/* arrange */
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									MinLength: len(*providedValue.String) - 1,
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
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Not: &model.StringConstraints{
										Pattern: "^.*$",
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
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Not: &model.StringConstraints{
										Pattern: "^$",
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
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									OneOf: []*model.StringConstraints{
										{
											Pattern: "^$",
										},
										{
											Pattern: *providedValue.String,
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
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									OneOf: []*model.StringConstraints{
										{
											Pattern: "^$",
										},
										{
											Enum: []string{"dummyEnumItem"},
										},
									},
								},
							},
						}

						expectedErrors := []error{
							errors.New("Must validate one and only one schema (oneOf)"),
							fmt.Errorf(
								`Does not match pattern '%v'`,
								providedParam.String.Constraints.OneOf[0].Pattern,
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
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									OneOf: []*model.StringConstraints{
										{
											Pattern: "^.*$",
										},
										{
											Enum: []string{*providedValue.String},
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
			Context("Pattern constraint", func() {
				Context("value doesn't match Pattern", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Pattern: "^$",
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf(
								"Does not match pattern '%v'",
								providedParam.String.Constraints.Pattern,
							),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value matches Pattern", func() {

					It("should return no errors", func() {

						/* arrange */
						providedValueString := "dummyValue"
						providedValue := &model.Value{
							String: &providedValueString,
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Pattern: ".$",
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
		Context("& nil value.String", func() {
			Context("and non empty Default", func() {
				Context("AllOf constraint", func() {
					Context("default meets all AllOf constraints", func() {

						It("returns no errors", func() {

							/* arrange */
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										AllOf: []*model.StringConstraints{
											{
												Pattern: "^.*$",
											},
											{
												Pattern: providedDefault,
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
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										AllOf: []*model.StringConstraints{
											{
												Pattern: "^$",
											},
											{
												Pattern: providedDefault,
											},
										},
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									`Does not match pattern '%v'`,
									providedParam.String.Constraints.AllOf[0].Pattern,
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
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										AnyOf: []*model.StringConstraints{
											{
												Pattern: "^.*$",
											},
											{
												Pattern: providedDefault,
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
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										AnyOf: []*model.StringConstraints{
											{
												Pattern: "^$",
											},
											{
												Enum: []string{"dummyEnumItem"},
											},
										},
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{
								errors.New("Must validate at least one schema (anyOf)"),
								fmt.Errorf(
									`Does not match pattern '%v'`,
									providedParam.String.Constraints.AnyOf[0].Pattern,
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
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Enum: []string{providedDefault},
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
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Enum: []string{"dummyEnumItem"},
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									`must be one of the following: "%v"`,
									providedParam.String.Constraints.Enum[0],
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
					Context("date-time", func() {
						Context("value doesn't match Format", func() {

							It("should return expected errors", func() {

								/* arrange */
								providedValue := &model.Value{}
								providedDefault := "notDateTime"
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "date-time",
										},
										Default: &providedDefault,
									},
								}

								expectedErrors := []error{
									fmt.Errorf(
										"Does not match format '%v'",
										providedParam.String.Constraints.Format,
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
								providedValue := &model.Value{}
								providedDefault := "0000-01-01T00:00:01.0Z"
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "date-time",
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
					Context("docker-image-ref", func() {
						Context("value doesn't match Format", func() {

							It("should return expected errors", func() {

								/* arrange */
								providedValue := &model.Value{}
								providedDefault := "$notADockerImageRef"
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "docker-image-ref",
										},
										Default: &providedDefault,
									},
								}

								expectedErrors := []error{
									fmt.Errorf(
										"Does not match format '%v'",
										providedParam.String.Constraints.Format,
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
								providedValue := &model.Value{}
								providedDefault := "dummy-registry.com/dummy-namespace/dummy-repo:dummy-tag"
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "docker-image-ref",
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
					Context("email", func() {
						Context("value doesn't match Format", func() {

							It("should return expected errors", func() {

								/* arrange */
								providedValue := &model.Value{}
								providedDefault := "notEmail"
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "email",
										},
										Default: &providedDefault,
									},
								}

								expectedErrors := []error{
									fmt.Errorf(
										"Does not match format '%v'",
										providedParam.String.Constraints.Format,
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
								providedValue := &model.Value{}
								providedDefault := "dummy-email@dummy-domain.com"
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "email",
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
					Context("hostname", func() {
						Context("value doesn't match Format", func() {

							It("should return expected errors", func() {

								/* arrange */
								providedValue := &model.Value{}
								providedDefault := "$notAHostname$"
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "hostname",
										},
										Default: &providedDefault,
									},
								}

								expectedErrors := []error{
									fmt.Errorf(
										"Does not match format '%v'",
										providedParam.String.Constraints.Format,
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
								providedValue := &model.Value{}
								providedDefault := "dummy.com"
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "hostname",
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
					Context("ipv4", func() {
						Context("value doesn't match Format", func() {

							It("should return expected errors", func() {

								/* arrange */
								providedValue := &model.Value{}
								providedDefault := "notAnIpV4"
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "ipv4",
										},
										Default: &providedDefault,
									},
								}

								expectedErrors := []error{
									fmt.Errorf(
										"Does not match format '%v'",
										providedParam.String.Constraints.Format,
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
								providedValue := &model.Value{}
								providedDefault := "0.0.0.0"
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "ipv4",
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
					Context("ipv6", func() {
						Context("value doesn't match Format", func() {

							It("should return expected errors", func() {

								/* arrange */
								providedValue := &model.Value{}
								providedDefault := "notAnIpV6"
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "ipv6",
										},
										Default: &providedDefault,
									},
								}

								expectedErrors := []error{
									fmt.Errorf(
										"Does not match format '%v'",
										providedParam.String.Constraints.Format,
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
								providedValue := &model.Value{}
								providedDefault := "0000:0000:0000:0000:0000:0000:0000:0000"
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "ipv6",
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
					Context("semver", func() {
						Context("value doesn't match Format", func() {

							It("should return expected errors", func() {

								/* arrange */
								providedValue := &model.Value{}
								providedDefault := "$notASemver$"
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "semver",
										},
										Default: &providedDefault,
									},
								}

								expectedErrors := []error{
									fmt.Errorf(
										"Does not match format '%v'",
										providedParam.String.Constraints.Format,
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
								providedValue := &model.Value{}
								providedDefault := "1.1.1"
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "semver",
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
					Context("uri", func() {
						Context("value doesn't match Format", func() {

							It("should return expected errors", func() {

								/* arrange */
								providedValue := &model.Value{}
								providedDefault := "notUri"
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "uri",
										},
										Default: &providedDefault,
									},
								}

								expectedErrors := []error{
									fmt.Errorf(
										"Does not match format '%v'",
										providedParam.String.Constraints.Format,
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
								providedValue := &model.Value{}
								providedDefault := "https://dummyuri.com:8080/somepath?somequery#somefragment"
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "uri",
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
				})
				Context("MaxLength constraint", func() {
					Context("default length == MaxLength", func() {

						It("returns no errors", func() {

							/* arrange */
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										MaxLength: len(providedDefault),
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
					Context("default length < MaxLength", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										MaxLength: len(providedDefault) - 1,
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"String length must be less than or equal to %v",
									providedParam.String.Constraints.MaxLength,
								),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default length > MaxLength", func() {

						It("returns no errors", func() {

							/* arrange */
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										MaxLength: len(providedDefault) + 1,
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
				Context("MinLength constraint", func() {
					Context("default length == MinLength", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										MinLength: len(providedDefault),
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
					Context("default length < MinLength", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										MinLength: len(providedDefault) + 1,
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"String length must be greater than or equal to %v",
									providedParam.String.Constraints.MinLength,
								),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default length > MinLength", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										MinLength: len(providedDefault) - 1,
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
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Not: &model.StringConstraints{
											Pattern: "^.*$",
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
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Not: &model.StringConstraints{
											Pattern: "^$",
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
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										OneOf: []*model.StringConstraints{
											{
												Pattern: "^$",
											},
											{
												Pattern: providedDefault,
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
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										OneOf: []*model.StringConstraints{
											{
												Pattern: "^$",
											},
											{
												Enum: []string{"dummyEnumItem"},
											},
										},
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{
								errors.New("Must validate one and only one schema (oneOf)"),
								fmt.Errorf(
									`Does not match pattern '%v'`,
									providedParam.String.Constraints.OneOf[0].Pattern,
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
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										OneOf: []*model.StringConstraints{
											{
												Pattern: "^.*$",
											},
											{
												Enum: []string{providedDefault},
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
				Context("Pattern constraint", func() {
					Context("default doesn't match Pattern", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Pattern: "^$",
									},
									Default: &providedDefault,
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"Does not match pattern '%v'",
									providedParam.String.Constraints.Pattern,
								),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default matches Pattern", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValue := &model.Value{}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Pattern: ".$",
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
			})
		})
	})
})
