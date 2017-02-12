package validate

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Describe("Param", func() {
	objectUnderTest := New()
	Context("invoked w/ nil param", func() {
		It("should panic", func() {
			/* arrange/act/assert */
			Expect(
				func() {
					objectUnderTest.Param(&model.Data{}, nil)
				},
			).To(Panic())
		})
	})
	Context("invoked w/ non-nil param.String", func() {
		Context("& non-empty value.String", func() {
			Context("AllOf constraint", func() {
				Context("value meets all AllOf constraints", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									AllOf: []*model.StringConstraints{
										{
											Pattern: "^.*$",
										},
										{
											Pattern: providedValue.String,
										},
									},
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value doesn't meet all AllOf constraints", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									AllOf: []*model.StringConstraints{
										{
											Pattern: "^$",
										},
										{
											Pattern: providedValue.String,
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
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
			Context("AnyOf constraint", func() {
				Context("value meets an AnyOf constraint", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									AnyOf: []*model.StringConstraints{
										{
											Pattern: "^.*$",
										},
										{
											Pattern: providedValue.String,
										},
									},
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value doesn't meet an AnyOf constraint", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
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
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
			Context("Enum constraint", func() {
				Context("value in enum", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Enum: []string{providedValue.String},
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value not in enum", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
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
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

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
							providedValue := &model.Data{
								String: "notDateTime",
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
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("value matches Format", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "0000-01-01T00:00:01.0Z",
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
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("docker-image-ref", func() {
					Context("value doesn't match Format", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "$notADockerImageRef",
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
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("value matches Format", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "dummy-registry.com/dummy-namespace/dummy-repo:dummy-tag",
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
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("email", func() {
					Context("value doesn't match Format", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "notEmail",
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
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("value matches Format", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "dummy-email@dummy-domain.com",
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
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("hostname", func() {
					Context("value doesn't match Format", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "$notAHostname$",
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
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("value matches Format", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "dummy.com",
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
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("ipv4", func() {
					Context("value doesn't match Format", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "notAnIpV4",
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
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("value matches Format", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "0.0.0.0",
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
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("ipv6", func() {
					Context("value doesn't match Format", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "notAnIpV6",
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
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("value matches Format", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "0000:0000:0000:0000:0000:0000:0000:0000",
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
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("uri", func() {
					Context("value doesn't match Format", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "notUri",
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
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("value matches Format", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "https://dummyuri.com:8080/somepath?somequery#somefragment",
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
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

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
						providedValue := &model.Data{
							String: "dummyValue",
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									MaxLength: len(providedValue.String),
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value length > MaxLength", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									MaxLength: len(providedValue.String) - 1,
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
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value length < MaxLength", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									MaxLength: len(providedValue.String) + 1,
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
			Context("MinLength constraint", func() {
				Context("value length == MinLength", func() {

					It("should return no errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									MinLength: len(providedValue.String),
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value length < MinLength", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									MinLength: len(providedValue.String) + 1,
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
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value length > MinLength", func() {

					It("should return no errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									MinLength: len(providedValue.String) - 1,
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
			Context("Not constraint", func() {
				Context("value matches", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
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
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value doesn't match", func() {

					It("should return no errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
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
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
			Context("OneOf constraint", func() {
				Context("value meets a single OneOf constraint", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									OneOf: []*model.StringConstraints{
										{
											Pattern: "^$",
										},
										{
											Pattern: providedValue.String,
										},
									},
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value meets no OneOf constraints", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
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
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value meets multiple OneOf constraints", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
						}
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									OneOf: []*model.StringConstraints{
										{
											Pattern: "^.*$",
										},
										{
											Enum: []string{providedValue.String},
										},
									},
								},
							},
						}

						expectedErrors := []error{
							errors.New("Must validate one and only one schema (oneOf)"),
						}

						/* act */
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
			Context("Pattern constraint", func() {
				Context("value doesn't match Pattern", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
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
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value matches Pattern", func() {

					It("should return no errors", func() {

						/* arrange */
						providedValue := &model.Data{
							String: "dummyValue",
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
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
		})
		Context("& empty value.String", func() {
			Context("and non empty Default", func() {
				Context("AllOf constraint", func() {
					Context("default meets all AllOf constraints", func() {

						It("returns no errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
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
									Default: providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("value doesn't meet all AllOf constraints", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
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
									Default: providedDefault,
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
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("AnyOf constraint", func() {
					Context("default meets an AnyOf constraint", func() {

						It("returns no errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
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
									Default: providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("value doesn't meet an AnyOf constraint", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
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
									Default: providedDefault,
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
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("Enum constraint", func() {
					Context("default in enum", func() {

						It("returns no errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Enum: []string{providedDefault},
									},
									Default: providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default not in enum", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Enum: []string{"dummyEnumItem"},
									},
									Default: providedDefault,
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									`must be one of the following: "%v"`,
									providedParam.String.Constraints.Enum[0],
								),
							}

							/* act */
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

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
								providedValue := &model.Data{
									String: "",
								}
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "date-time",
										},
										Default: "notDateTime",
									},
								}

								expectedErrors := []error{
									fmt.Errorf(
										"Does not match format '%v'",
										providedParam.String.Constraints.Format,
									),
								}

								/* act */
								actualErrors := objectUnderTest.Param(providedValue, providedParam)

								/* assert */
								Expect(actualErrors).To(Equal(expectedErrors))

							})
						})
						Context("value matches Format", func() {

							It("should return no errors", func() {

								/* arrange */
								providedValue := &model.Data{
									String: "",
								}
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "date-time",
										},
										Default: "0000-01-01T00:00:01.0Z",
									},
								}

								expectedErrors := []error{}

								/* act */
								actualErrors := objectUnderTest.Param(providedValue, providedParam)

								/* assert */
								Expect(actualErrors).To(Equal(expectedErrors))

							})
						})
					})
					Context("docker-image-ref", func() {
						Context("value doesn't match Format", func() {

							It("should return expected errors", func() {

								/* arrange */
								providedValue := &model.Data{
									String: "",
								}
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "docker-image-ref",
										},
										Default: "$notADockerImageRef",
									},
								}

								expectedErrors := []error{
									fmt.Errorf(
										"Does not match format '%v'",
										providedParam.String.Constraints.Format,
									),
								}

								/* act */
								actualErrors := objectUnderTest.Param(providedValue, providedParam)

								/* assert */
								Expect(actualErrors).To(Equal(expectedErrors))

							})
						})
						Context("value matches Format", func() {

							It("should return no errors", func() {

								/* arrange */
								providedValue := &model.Data{
									String: "",
								}
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "docker-image-ref",
										},
										Default: "dummy-registry.com/dummy-namespace/dummy-repo:dummy-tag",
									},
								}

								expectedErrors := []error{}

								/* act */
								actualErrors := objectUnderTest.Param(providedValue, providedParam)

								/* assert */
								Expect(actualErrors).To(Equal(expectedErrors))

							})
						})
					})
					Context("email", func() {
						Context("value doesn't match Format", func() {

							It("should return expected errors", func() {

								/* arrange */
								providedValue := &model.Data{
									String: "",
								}
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "email",
										},
										Default: "notEmail",
									},
								}

								expectedErrors := []error{
									fmt.Errorf(
										"Does not match format '%v'",
										providedParam.String.Constraints.Format,
									),
								}

								/* act */
								actualErrors := objectUnderTest.Param(providedValue, providedParam)

								/* assert */
								Expect(actualErrors).To(Equal(expectedErrors))

							})
						})
						Context("value matches Format", func() {

							It("should return no errors", func() {

								/* arrange */
								providedValue := &model.Data{
									String: "",
								}
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "email",
										},
										Default: "dummy-email@dummy-domain.com",
									},
								}

								expectedErrors := []error{}

								/* act */
								actualErrors := objectUnderTest.Param(providedValue, providedParam)

								/* assert */
								Expect(actualErrors).To(Equal(expectedErrors))

							})
						})
					})
					Context("hostname", func() {
						Context("value doesn't match Format", func() {

							It("should return expected errors", func() {

								/* arrange */
								providedValue := &model.Data{
									String: "",
								}
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "hostname",
										},
										Default: "$notAHostname$",
									},
								}

								expectedErrors := []error{
									fmt.Errorf(
										"Does not match format '%v'",
										providedParam.String.Constraints.Format,
									),
								}

								/* act */
								actualErrors := objectUnderTest.Param(providedValue, providedParam)

								/* assert */
								Expect(actualErrors).To(Equal(expectedErrors))

							})
						})
						Context("value matches Format", func() {

							It("should return no errors", func() {

								/* arrange */
								providedValue := &model.Data{
									String: "",
								}
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "hostname",
										},
										Default: "dummy.com",
									},
								}

								expectedErrors := []error{}

								/* act */
								actualErrors := objectUnderTest.Param(providedValue, providedParam)

								/* assert */
								Expect(actualErrors).To(Equal(expectedErrors))

							})
						})
					})
					Context("ipv4", func() {
						Context("value doesn't match Format", func() {

							It("should return expected errors", func() {

								/* arrange */
								providedValue := &model.Data{
									String: "",
								}
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "ipv4",
										},
										Default: "notAnIpV4",
									},
								}

								expectedErrors := []error{
									fmt.Errorf(
										"Does not match format '%v'",
										providedParam.String.Constraints.Format,
									),
								}

								/* act */
								actualErrors := objectUnderTest.Param(providedValue, providedParam)

								/* assert */
								Expect(actualErrors).To(Equal(expectedErrors))

							})
						})
						Context("value matches Format", func() {

							It("should return no errors", func() {

								/* arrange */
								providedValue := &model.Data{
									String: "",
								}
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "ipv4",
										},
										Default: "0.0.0.0",
									},
								}

								expectedErrors := []error{}

								/* act */
								actualErrors := objectUnderTest.Param(providedValue, providedParam)

								/* assert */
								Expect(actualErrors).To(Equal(expectedErrors))

							})
						})
					})
					Context("ipv6", func() {
						Context("value doesn't match Format", func() {

							It("should return expected errors", func() {

								/* arrange */
								providedValue := &model.Data{
									String: "",
								}
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "ipv6",
										},
										Default: "notAnIpV6",
									},
								}

								expectedErrors := []error{
									fmt.Errorf(
										"Does not match format '%v'",
										providedParam.String.Constraints.Format,
									),
								}

								/* act */
								actualErrors := objectUnderTest.Param(providedValue, providedParam)

								/* assert */
								Expect(actualErrors).To(Equal(expectedErrors))

							})
						})
						Context("value matches Format", func() {

							It("should return no errors", func() {

								/* arrange */
								providedValue := &model.Data{
									String: "",
								}
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "ipv6",
										},
										Default: "0000:0000:0000:0000:0000:0000:0000:0000",
									},
								}

								expectedErrors := []error{}

								/* act */
								actualErrors := objectUnderTest.Param(providedValue, providedParam)

								/* assert */
								Expect(actualErrors).To(Equal(expectedErrors))

							})
						})
					})
					Context("uri", func() {
						Context("value doesn't match Format", func() {

							It("should return expected errors", func() {

								/* arrange */
								providedValue := &model.Data{
									String: "",
								}
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "uri",
										},
										Default: "notUri",
									},
								}

								expectedErrors := []error{
									fmt.Errorf(
										"Does not match format '%v'",
										providedParam.String.Constraints.Format,
									),
								}

								/* act */
								actualErrors := objectUnderTest.Param(providedValue, providedParam)

								/* assert */
								Expect(actualErrors).To(Equal(expectedErrors))

							})
						})
						Context("value matches Format", func() {

							It("should return no errors", func() {

								/* arrange */
								providedValue := &model.Data{
									String: "",
								}
								providedParam := &model.Param{
									String: &model.StringParam{
										Constraints: &model.StringConstraints{
											Format: "uri",
										},
										Default: "https://dummyuri.com:8080/somepath?somequery#somefragment",
									},
								}

								expectedErrors := []error{}

								/* act */
								actualErrors := objectUnderTest.Param(providedValue, providedParam)

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
							providedValue := &model.Data{
								String: "",
							}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										MaxLength: len(providedDefault),
									},
									Default: providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default length < MaxLength", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										MaxLength: len(providedDefault) - 1,
									},
									Default: providedDefault,
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"String length must be less than or equal to %v",
									providedParam.String.Constraints.MaxLength,
								),
							}

							/* act */
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default length > MaxLength", func() {

						It("returns no errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										MaxLength: len(providedDefault) + 1,
									},
									Default: providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("MinLength constraint", func() {
					Context("default length == MinLength", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										MinLength: len(providedDefault),
									},
									Default: providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default length < MinLength", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										MinLength: len(providedDefault) + 1,
									},
									Default: providedDefault,
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"String length must be greater than or equal to %v",
									providedParam.String.Constraints.MinLength,
								),
							}

							/* act */
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default length > MinLength", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										MinLength: len(providedDefault) - 1,
									},
									Default: providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("Not constraint", func() {
					Context("default matches", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Not: &model.StringConstraints{
											Pattern: "^.*$",
										},
									},
									Default: providedDefault,
								},
							}

							expectedErrors := []error{
								errors.New("Must not validate the schema (not)"),
							}

							/* act */
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default doesn't match", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Not: &model.StringConstraints{
											Pattern: "^$",
										},
									},
									Default: providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("OneOf constraint", func() {
					Context("default meets a single OneOf constraint", func() {

						It("returns no errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
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
									Default: providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default meets no OneOf constraints", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
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
									Default: providedDefault,
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
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default meets multiple OneOf constraints", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
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
									Default: providedDefault,
								},
							}

							expectedErrors := []error{
								errors.New("Must validate one and only one schema (oneOf)"),
							}

							/* act */
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("Pattern constraint", func() {
					Context("default doesn't match Pattern", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Pattern: "^$",
									},
									Default: providedDefault,
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"Does not match pattern '%v'",
									providedParam.String.Constraints.Pattern,
								),
							}

							/* act */
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default matches Pattern", func() {

						It("should return no errors", func() {

							/* arrange */
							providedValue := &model.Data{
								String: "",
							}
							providedDefault := "dummyDefault"
							providedParam := &model.Param{
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Pattern: ".$",
									},
									Default: providedDefault,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Param(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
			})
		})
		Context("& nil value", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedParam := &model.Param{
					String: &model.StringParam{},
				}

				expectedErrors := []error{
					errors.New("String required"),
				}

				/* act */
				actualErrors := objectUnderTest.Param(nil, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
	})
	Context("invoked w/ non-nil param.Socket", func() {
		Context("& non-empty value.Socket", func() {
			It("should return no errors", func() {

				/* arrange */
				providedValue := &model.Data{
					Socket: "dummyValue",
				}
				providedParam := &model.Param{
					Socket: &model.SocketParam{},
				}

				expectedErrors := []error{}

				/* act */
				actualErrors := objectUnderTest.Param(providedValue, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
		Context("& empty value.Socket", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedValue := &model.Data{}
				providedParam := &model.Param{
					Socket: &model.SocketParam{},
				}

				expectedErrors := []error{
					errors.New("Socket required"),
				}

				/* act */
				actualErrors := objectUnderTest.Param(providedValue, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
		Context("& nil value", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedParam := &model.Param{
					Socket: &model.SocketParam{},
				}

				expectedErrors := []error{
					errors.New("Socket required"),
				}

				/* act */
				actualErrors := objectUnderTest.Param(nil, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
	})

})
