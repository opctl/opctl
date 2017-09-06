package string

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
				errors.New("string required"),
			}

			objectUnderTest := newValidator()

			/* act */
			actualErrors := objectUnderTest.Validate(
				nil,
				&model.StringConstraints{},
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
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						AllOf: []*model.StringConstraints{
							{
								Pattern: "^.*$",
							},
							{
								Pattern: providedValue,
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
			Context("value doesn't meet all AllOf constraints", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValue := "dummyValue==\""
					providedConstraints := &model.StringConstraints{
						AllOf: []*model.StringConstraints{
							{
								Pattern: "^$",
							},
							{
								Pattern: providedValue,
							},
						},
					}

					expectedErrors := []error{
						fmt.Errorf(
							`Does not match pattern '%v'`,
							providedConstraints.AllOf[0].Pattern,
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
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						AnyOf: []*model.StringConstraints{
							{
								Pattern: "^.*$",
							},
							{
								Pattern: providedValue,
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
			Context("value doesn't meet an AnyOf constraint", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						AnyOf: []*model.StringConstraints{
							{
								Pattern: "^$",
							},
							{
								Enum: []string{"dummyEnumItem"},
							},
						},
					}

					expectedErrors := []error{
						errors.New("Must validate at least one schema (anyOf)"),
						fmt.Errorf(
							`Does not match pattern '%v'`,
							providedConstraints.AnyOf[0].Pattern,
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
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						Enum: []string{providedValue},
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
			Context("value not in enum", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						Enum: []string{"dummyEnumItem"},
					}

					expectedErrors := []error{
						fmt.Errorf(
							`must be one of the following: "%v"`,
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
			Context("date-time", func() {
				Context("value doesn't match Format", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValue := "notDateTime"
						providedConstraints := &model.StringConstraints{
							Format: "date-time",
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
						providedValue := "0000-01-01T00:00:01.0Z"
						providedConstraints := &model.StringConstraints{
							Format: "date-time",
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
			Context("docker-image-ref", func() {
				Context("value doesn't match Format", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValue := "$notADockerImageRef"
						providedConstraints := &model.StringConstraints{
							Format: "docker-image-ref",
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
						providedValue := "dummy-registry.com/dummy-namespace/dummy-repo:dummy-tag"
						providedConstraints := &model.StringConstraints{
							Format: "docker-image-ref",
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
			Context("email", func() {
				Context("value doesn't match Format", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValue := "notEmail"
						providedConstraints := &model.StringConstraints{
							Format: "email",
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
						providedValue := "dummy-email@dummy-domain.com"
						providedConstraints := &model.StringConstraints{
							Format: "email",
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
			Context("hostname", func() {
				Context("value doesn't match Format", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValue := "$notAHostname$"
						providedConstraints := &model.StringConstraints{
							Format: "hostname",
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
						providedValue := "dummy.com"
						providedConstraints := &model.StringConstraints{
							Format: "hostname",
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
			Context("ipv4", func() {
				Context("value doesn't match Format", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValue := "notAnIpV4"
						providedConstraints := &model.StringConstraints{
							Format: "ipv4",
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
						providedValue := "0.0.0.0"
						providedConstraints := &model.StringConstraints{
							Format: "ipv4",
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
			Context("ipv6", func() {
				Context("value doesn't match Format", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValue := "notAnIpV6"
						providedConstraints := &model.StringConstraints{
							Format: "ipv6",
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
						providedValue := "0000:0000:0000:0000:0000:0000:0000:0000"
						providedConstraints := &model.StringConstraints{
							Format: "ipv6",
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
			Context("semver", func() {
				Context("value doesn't match Format", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValue := "$notASemver$"
						providedConstraints := &model.StringConstraints{
							Format: "semver",
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
						providedValue := "1.1.1"
						providedConstraints := &model.StringConstraints{
							Format: "semver",
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
			Context("uri", func() {
				Context("value doesn't match Format", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValue := "notUri"
						providedConstraints := &model.StringConstraints{
							Format: "uri",
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
						providedValue := "https://dummyuri.com:8080/somepath?somequery#somefragment"
						providedConstraints := &model.StringConstraints{
							Format: "uri",
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
		Context("MaxLength constraint", func() {
			Context("value length == MaxLength", func() {

				It("returns no errors", func() {

					/* arrange */
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						MaxLength: len(providedValue),
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
			Context("value length > MaxLength", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						MaxLength: len(providedValue) - 1,
					}

					expectedErrors := []error{
						fmt.Errorf(
							"String length must be less than or equal to %v",
							providedConstraints.MaxLength,
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
			Context("value length < MaxLength", func() {

				It("returns no errors", func() {

					/* arrange */
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						MaxLength: len(providedValue) + 1,
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
		Context("MinLength constraint", func() {
			Context("value length == MinLength", func() {

				It("should return no errors", func() {

					/* arrange */
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						MinLength: len(providedValue),
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
			Context("value length < MinLength", func() {

				It("should return expected errors", func() {

					/* arrange */
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						MinLength: len(providedValue) + 1,
					}

					expectedErrors := []error{
						fmt.Errorf(
							"String length must be greater than or equal to %v",
							providedConstraints.MinLength,
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
			Context("value length > MinLength", func() {

				It("should return no errors", func() {

					/* arrange */
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						MinLength: len(providedValue) - 1,
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
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						Not: &model.StringConstraints{
							Pattern: "^.*$",
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
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						Not: &model.StringConstraints{
							Pattern: "^$",
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
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						OneOf: []*model.StringConstraints{
							{
								Pattern: "^$",
							},
							{
								Pattern: providedValue,
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
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						OneOf: []*model.StringConstraints{
							{
								Pattern: "^$",
							},
							{
								Enum: []string{"dummyEnumItem"},
							},
						},
					}

					expectedErrors := []error{
						errors.New("Must validate one and only one schema (oneOf)"),
						fmt.Errorf(
							`Does not match pattern '%v'`,
							providedConstraints.OneOf[0].Pattern,
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
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						OneOf: []*model.StringConstraints{
							{
								Pattern: "^.*$",
							},
							{
								Enum: []string{providedValue},
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
		Context("Pattern constraint", func() {
			Context("value doesn't match Pattern", func() {

				It("should return expected errors", func() {

					/* arrange */
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						Pattern: "^$",
					}

					expectedErrors := []error{
						fmt.Errorf(
							"Does not match pattern '%v'",
							providedConstraints.Pattern,
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
			Context("value matches Pattern", func() {

				It("should return no errors", func() {

					/* arrange */
					providedValue := "dummyValue"
					providedConstraints := &model.StringConstraints{
						Pattern: ".$",
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
})
