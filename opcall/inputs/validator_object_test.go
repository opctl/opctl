package inputs

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"strings"
)

var _ = Context("Validate", func() {
	objectUnderTest := newValidator()
	Context("invoked w/ non-nil param.Object", func() {
		Context("& non-nil value.Object", func() {
			Context("AdditionalProperties constraint", func() {
				Context("value props don't match properties, or patternProperties", func() {
					Context("value props meet AdditionalProperties constraint", func() {

						It("returns no errors", func() {

							/* arrange */
							providedValueObject := map[string]interface{}{
								"dummyProp1Name": "dummyProp1Value",
							}
							providedValue := &model.Value{
								Object: providedValueObject,
							}
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										AdditionalProperties: &model.JSONSchema{
											MinLength: 2,
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
					Context("some value props meet AdditionalProperties constraint", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedValueObjectProp1Value := "dummyProp1Value"
							providedValueObject := map[string]interface{}{
								"dummyProp1Name": providedValueObjectProp1Value,
								"dummyProp2Name": "dummyProp2Value",
							}
							providedValue := &model.Value{
								Object: providedValueObject,
							}

							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										AdditionalProperties: &model.JSONSchema{
											Pattern: providedValueObjectProp1Value,
										},
									},
								},
							}

							expectedErrors := []error{
								fmt.Errorf("Does not match pattern '%v'", providedValueObjectProp1Value),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("no value props meet AdditionalProperties constraint", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedValueObject := map[string]interface{}{
								"dummyProp1Name": "dummyProp1Value",
								"dummyProp2Name": "dummyProp2Value",
							}
							providedValue := &model.Value{
								Object: providedValueObject,
							}

							pattern := "dummyPattern"
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										AdditionalProperties: &model.JSONSchema{
											Pattern: pattern,
										},
									},
								},
							}

							expectedErrors := []error{
								fmt.Errorf("Does not match pattern '%v'", pattern),
								fmt.Errorf("Does not match pattern '%v'", pattern),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(providedValue, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
			})
			Context("AllOf constraint", func() {
				Context("value meets all AllOf constraints", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValueObject := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									AllOf: []*model.ObjectConstraints{
										{
											MinProperties: 0,
										},
										{
											MaxProperties: 1,
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
						providedValueObject := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
							"dummyProp2Name": "dummyProp2Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									AllOf: []*model.ObjectConstraints{
										{
											MinProperties: 1,
										},
										{
											MaxProperties: 1,
										},
									},
								},
							},
						}

						expectedErrors := []error{
							errors.New(`Must have at most 1 properties`),
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
						providedValueObject := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									AnyOf: []*model.ObjectConstraints{
										{
											MinProperties: 1,
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
						providedValueObject := map[string]interface{}{}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									AnyOf: []*model.ObjectConstraints{
										{
											MinProperties: 1,
										},
									},
								},
							},
						}

						expectedErrors := []error{
							errors.New("Must validate at least one schema (anyOf)"),
							errors.New("Must have at least 1 properties"),
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
						providedValueObject := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									Enum: []map[string]interface{}{
										providedValueObject,
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
				Context("value not in enum", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValueObject := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									Enum: []map[string]interface{}{
										{
											"dummyName": "dummyValue",
										},
									},
								},
							},
						}

						expectedError, err := json.Marshal(providedParam.Object.Constraints.Enum[0])
						if nil != err {
							Fail(err.Error())
						}

						expectedErrors := []error{
							fmt.Errorf(
								`must be one of the following: %v`,
								string(expectedError),
							),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
			Context("MaxProperties constraint", func() {
				Context("value prop count == MaxProperties", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValueObject := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									MaxProperties: 1,
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
				Context("value prop count > MaxProperties", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValueObject := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
							"dummyProp2Name": "dummyProp2Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									MaxProperties: 1,
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf(
								"Must have at most %v properties",
								providedParam.Object.Constraints.MaxProperties,
							),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value prop count < MaxProperties", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValueObject := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									MaxProperties: 2,
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
			Context("MinProperties constraint", func() {
				Context("value prop count == MinProperties", func() {

					It("should return no errors", func() {

						/* arrange */
						providedValueObject := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									MinProperties: 1,
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
				Context("value prop count < MinProperties", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValueObject := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									MinProperties: 2,
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf(
								"Must have at least %v properties",
								providedParam.Object.Constraints.MinProperties,
							),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value prop count > MinProperties", func() {

					It("should return no errors", func() {

						/* arrange */
						providedValueObject := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
							"dummyProp2Name": "dummyProp2Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									MinProperties: 1,
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
						providedValueObject := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									Not: &model.ObjectConstraints{
										MinProperties: 1,
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
						providedValueObject := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									Not: &model.ObjectConstraints{
										MinProperties: 2,
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
						providedValueObject := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									OneOf: []*model.ObjectConstraints{
										{
											MinProperties: 1,
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
						providedValueObject := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									OneOf: []*model.ObjectConstraints{
										{
											MinProperties: 2,
										},
									},
								},
							},
						}

						expectedErrors := []error{
							errors.New("Must validate one and only one schema (oneOf)"),
							errors.New("Must have at least 2 properties"),
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
						providedValueObject := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
							"dummyProp2Name": "dummyProp2Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									OneOf: []*model.ObjectConstraints{
										{
											MinProperties: 2,
										},
										{
											MaxProperties: 2,
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
			Context("PatternProperties constraint", func() {
				Context("value props meet all Properties constraints", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValueObjectProp1Name := "dummyProp1Name"
						providedValueObject := map[string]interface{}{
							providedValueObjectProp1Name: "dummyProp1Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									PatternProperties: map[string]*model.JSONSchema{
										providedValueObjectProp1Name: {MinLength: 2},
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
				Context("value props meet some Properties constraints", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValueObjectProp1Name := "dummyProp1Name"
						providedValueObjectProp2Name := "dummyProp2Name"
						providedValueObject := map[string]interface{}{
							providedValueObjectProp1Name: "dummyProp1Value",
							providedValueObjectProp2Name: "dummyProp2Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}

						maxLength := 1
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									PatternProperties: map[string]*model.JSONSchema{
										providedValueObjectProp1Name: {MinLength: 2},
										providedValueObjectProp2Name: {MaxLength: maxLength},
									},
								},
							},
						}

						patterns, err := json.Marshal([]string{providedValueObjectProp1Name, providedValueObjectProp2Name})
						if err != nil {
							Fail(err.Error())
						}

						expectedErrors := []error{
							fmt.Errorf("String length must be less than or equal to %v", maxLength),
							fmt.Errorf("Property \"%v\" does not match pattern %v", providedValueObjectProp2Name, string(patterns)),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						// fix pattern order (maps in go not stable)
						OutOfOrderPatterns, err := json.Marshal([]string{providedValueObjectProp2Name, providedValueObjectProp1Name})
						if err != nil {
							Fail(err.Error())
						}
						for errIndex, err := range actualErrors {
							actualErrors[errIndex] = errors.New(
								strings.Replace(err.Error(), string(OutOfOrderPatterns), string(patterns), -1),
							)
						}
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value props meet no Properties constraints", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValueObjectProp1Name := "dummyProp1Name"
						providedValueObjectProp2Name := "dummyProp2Name"
						providedValueObject := map[string]interface{}{
							providedValueObjectProp1Name: "dummyProp1Value",
							providedValueObjectProp2Name: "dummyProp2Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						minLength := 100
						maxLength := 1
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									PatternProperties: map[string]*model.JSONSchema{
										providedValueObjectProp1Name: {MinLength: minLength},
										providedValueObjectProp2Name: {MaxLength: maxLength},
									},
								},
							},
						}

						patterns, err := json.Marshal([]string{providedValueObjectProp1Name, providedValueObjectProp2Name})
						if err != nil {
							Fail(err.Error())
						}

						expectedErrors := []error{
							fmt.Errorf("String length must be greater than or equal to %v", minLength),
							fmt.Errorf("Property \"%v\" does not match pattern %v", providedValueObjectProp1Name, string(patterns)),
							fmt.Errorf("String length must be less than or equal to %v", maxLength),
							fmt.Errorf("Property \"%v\" does not match pattern %v", providedValueObjectProp2Name, string(patterns)),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						// fix pattern order (maps in go not stable)
						OutOfOrderPatterns, err := json.Marshal([]string{providedValueObjectProp2Name, providedValueObjectProp1Name})
						if err != nil {
							Fail(err.Error())
						}
						for errIndex, err := range actualErrors {
							actualErrors[errIndex] = errors.New(
								strings.Replace(err.Error(), string(OutOfOrderPatterns), string(patterns), -1),
							)
						}
						Expect(actualErrors).To(ConsistOf(expectedErrors))

					})
				})
			})
			Context("Properties constraint", func() {
				Context("value props meet all Properties constraints", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValueObjectProp1Name := "dummyProp1Name"
						providedValueObject := map[string]interface{}{
							providedValueObjectProp1Name: "dummyProp1Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									Properties: map[string]*model.JSONSchema{
										providedValueObjectProp1Name: {MinLength: 2},
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
				Context("value props meet some Properties constraints", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValueObjectProp1Name := "dummyProp1Name"
						providedValueObjectProp2Name := "dummyProp2Name"
						providedValueObject := map[string]interface{}{
							providedValueObjectProp1Name: "dummyProp1Value",
							providedValueObjectProp2Name: "dummyProp2Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}

						maxLength := 1
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									Properties: map[string]*model.JSONSchema{
										providedValueObjectProp1Name: {MinLength: 2},
										providedValueObjectProp2Name: {MaxLength: maxLength},
									},
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf("String length must be less than or equal to %v", maxLength),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value props meet no Properties constraints", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValueObjectProp1Name := "dummyProp1Name"
						providedValueObjectProp2Name := "dummyProp2Name"
						providedValueObject := map[string]interface{}{
							providedValueObjectProp1Name: "dummyProp1Value",
							providedValueObjectProp2Name: "dummyProp2Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						minLength := 100
						maxLength := 1
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									Properties: map[string]*model.JSONSchema{
										providedValueObjectProp1Name: {MinLength: minLength},
										providedValueObjectProp2Name: {MaxLength: maxLength},
									},
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf("String length must be greater than or equal to %v", minLength),
							fmt.Errorf("String length must be less than or equal to %v", maxLength),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(ConsistOf(expectedErrors))

					})
				})
			})
			Context("Required constraint", func() {
				Context("value contains all required props", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValueObjectProp1Name := "dummyProp1Name"
						providedValueObject := map[string]interface{}{
							providedValueObjectProp1Name: "dummyProp1Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									Required: []string{
										providedValueObjectProp1Name,
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
				Context("value contains no required props", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValueObject := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}

						missingPropName := "missingPropName"
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									Required: []string{
										missingPropName,
									},
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf(
								"%v is required",
								missingPropName,
							),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value contains some required props", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValueObjectProp1Name := "dummyProp1Name"
						providedValueObject := map[string]interface{}{
							providedValueObjectProp1Name: "dummyProp1Value",
							"dummyProp2Name":             "dummyProp2Value",
						}
						providedValue := &model.Value{
							Object: providedValueObject,
						}

						missingPropName := "missingPropName"
						providedParam := &model.Param{
							Object: &model.ObjectParam{
								Constraints: &model.ObjectConstraints{
									Required: []string{
										providedValueObjectProp1Name,
										missingPropName,
									},
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf(
								"%v is required",
								missingPropName,
							),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
		})
		Context("& nil value.Object", func() {
			Context("and non-nil Default", func() {
				Context("AdditionalProperties constraint", func() {
					Context("default props don't match properties, or patternProperties", func() {
						Context("default props meet AdditionalProperties constraint", func() {

							It("returns no errors", func() {

								/* arrange */
								providedParam := &model.Param{
									Object: &model.ObjectParam{
										Constraints: &model.ObjectConstraints{
											AdditionalProperties: &model.JSONSchema{
												MinLength: 2,
											},
										},
										Default: map[string]interface{}{
											"dummyProp1Name": "dummyProp1Value",
										},
									},
								}

								expectedErrors := []error{}

								/* act */
								actualErrors := objectUnderTest.Validate(nil, providedParam)

								/* assert */
								Expect(actualErrors).To(Equal(expectedErrors))

							})
						})
						Context("some default props meet AdditionalProperties constraint", func() {

							It("returns expected errors", func() {

								/* arrange */
								defaultObjectProp1Value := "dummyProp1Value"

								providedParam := &model.Param{
									Object: &model.ObjectParam{
										Constraints: &model.ObjectConstraints{
											AdditionalProperties: &model.JSONSchema{
												Pattern: defaultObjectProp1Value,
											},
										},
										Default: map[string]interface{}{
											"dummyProp1Name": defaultObjectProp1Value,
											"dummyProp2Name": "dummyProp2Value",
										},
									},
								}

								expectedErrors := []error{
									fmt.Errorf("Does not match pattern '%v'", defaultObjectProp1Value),
								}

								/* act */
								actualErrors := objectUnderTest.Validate(nil, providedParam)

								/* assert */
								Expect(actualErrors).To(Equal(expectedErrors))

							})
						})
						Context("no default props meet AdditionalProperties constraint", func() {

							It("returns expected errors", func() {

								/* arrange */
								pattern := "dummyPattern"
								providedParam := &model.Param{
									Object: &model.ObjectParam{
										Constraints: &model.ObjectConstraints{
											AdditionalProperties: &model.JSONSchema{
												Pattern: pattern,
											},
										},
										Default: map[string]interface{}{
											"dummyProp1Name": "dummyProp1Value",
											"dummyProp2Name": "dummyProp2Value",
										},
									},
								}

								expectedErrors := []error{
									fmt.Errorf("Does not match pattern '%v'", pattern),
									fmt.Errorf("Does not match pattern '%v'", pattern),
								}

								/* act */
								actualErrors := objectUnderTest.Validate(nil, providedParam)

								/* assert */
								Expect(actualErrors).To(Equal(expectedErrors))

							})
						})
					})
				})
				Context("AllOf constraint", func() {
					Context("default meets all AllOf constraints", func() {

						It("returns no errors", func() {

							/* arrange */
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										AllOf: []*model.ObjectConstraints{
											{
												MinProperties: 0,
											},
											{
												MaxProperties: 1,
											},
										},
									},
									Default: map[string]interface{}{
										"dummyProp1Name": "dummyProp1Value",
									},
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default doesn't meet all AllOf constraints", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										AllOf: []*model.ObjectConstraints{
											{
												MinProperties: 1,
											},
											{
												MaxProperties: 1,
											},
										},
									},
									Default: map[string]interface{}{
										"dummyProp1Name": "dummyProp1Value",
										"dummyProp2Name": "dummyProp2Value",
									},
								},
							}

							expectedErrors := []error{
								errors.New(`Must have at most 1 properties`),
								errors.New("Must validate all the schemas (allOf)"),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("AnyOf constraint", func() {
					Context("default meets an AnyOf constraint", func() {

						It("returns no errors", func() {

							/* arrange */
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										AnyOf: []*model.ObjectConstraints{
											{
												MinProperties: 1,
											},
										},
									},
									Default: map[string]interface{}{
										"dummyProp1Name": "dummyProp1Value",
									},
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default doesn't meet an AnyOf constraint", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										AnyOf: []*model.ObjectConstraints{
											{
												MinProperties: 1,
											},
										},
									},
									Default: map[string]interface{}{},
								},
							}

							expectedErrors := []error{
								errors.New("Must validate at least one schema (anyOf)"),
								errors.New("Must have at least 1 properties"),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("Enum constraint", func() {
					Context("default in enum", func() {

						It("returns no errors", func() {

							/* arrange */
							defaultObject := map[string]interface{}{
								"dummyProp1Name": "dummyProp1Value",
							}
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										Enum: []map[string]interface{}{
											defaultObject,
										},
									},
									Default: defaultObject,
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default not in enum", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										Enum: []map[string]interface{}{
											{
												"dummyName": "dummyValue",
											},
										},
									},
									Default: map[string]interface{}{
										"dummyProp1Name": "dummyProp1Value",
									},
								},
							}

							expectedError, err := json.Marshal(providedParam.Object.Constraints.Enum[0])
							if nil != err {
								Fail(err.Error())
							}

							expectedErrors := []error{
								fmt.Errorf(
									`must be one of the following: %v`,
									string(expectedError),
								),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("MaxProperties constraint", func() {
					Context("default prop count == MaxProperties", func() {

						It("returns no errors", func() {

							/* arrange */
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										MaxProperties: 1,
									},
									Default: map[string]interface{}{
										"dummyProp1Name": "dummyProp1Value",
									},
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default prop count > MaxProperties", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										MaxProperties: 1,
									},
									Default: map[string]interface{}{
										"dummyProp1Name": "dummyProp1Value",
										"dummyProp2Name": "dummyProp2Value",
									},
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"Must have at most %v properties",
									providedParam.Object.Constraints.MaxProperties,
								),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default prop count < MaxProperties", func() {

						It("returns no errors", func() {

							/* arrange */
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										MaxProperties: 2,
									},
									Default: map[string]interface{}{
										"dummyProp1Name": "dummyProp1Value",
									},
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("MinProperties constraint", func() {
					Context("default prop count == MinProperties", func() {

						It("should return no errors", func() {

							/* arrange */
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										MinProperties: 1,
									},
									Default: map[string]interface{}{
										"dummyProp1Name": "dummyProp1Value",
									},
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default prop count < MinProperties", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										MinProperties: 2,
									},
									Default: map[string]interface{}{
										"dummyProp1Name": "dummyProp1Value",
									},
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"Must have at least %v properties",
									providedParam.Object.Constraints.MinProperties,
								),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default prop count > MinProperties", func() {

						It("should return no errors", func() {

							/* arrange */
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										MinProperties: 1,
									},
									Default: map[string]interface{}{
										"dummyProp1Name": "dummyProp1Value",
										"dummyProp2Name": "dummyProp2Value",
									},
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("Not constraint", func() {
					Context("default matches", func() {

						It("should return expected errors", func() {

							/* arrange */
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										Not: &model.ObjectConstraints{
											MinProperties: 1,
										},
									},
									Default: map[string]interface{}{
										"dummyProp1Name": "dummyProp1Value",
									},
								},
							}

							expectedErrors := []error{
								errors.New("Must not validate the schema (not)"),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default doesn't match", func() {

						It("should return no errors", func() {

							/* arrange */
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										Not: &model.ObjectConstraints{
											MinProperties: 2,
										},
									},
									Default: map[string]interface{}{
										"dummyProp1Name": "dummyProp1Value",
									},
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("OneOf constraint", func() {
					Context("default meets a single OneOf constraint", func() {

						It("returns no errors", func() {

							/* arrange */
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										OneOf: []*model.ObjectConstraints{
											{
												MinProperties: 1,
											},
										},
									},
									Default: map[string]interface{}{
										"dummyProp1Name": "dummyProp1Value",
									},
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default meets no OneOf constraints", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										OneOf: []*model.ObjectConstraints{
											{
												MinProperties: 2,
											},
										},
									},
									Default: map[string]interface{}{
										"dummyProp1Name": "dummyProp1Value",
									},
								},
							}

							expectedErrors := []error{
								errors.New("Must validate one and only one schema (oneOf)"),
								errors.New("Must have at least 2 properties"),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default meets multiple OneOf constraints", func() {

						It("returns expected errors", func() {

							/* arrange */
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										OneOf: []*model.ObjectConstraints{
											{
												MinProperties: 2,
											},
											{
												MaxProperties: 2,
											},
										},
									},
									Default: map[string]interface{}{
										"dummyProp1Name": "dummyProp1Value",
										"dummyProp2Name": "dummyProp2Value",
									},
								},
							}

							expectedErrors := []error{
								errors.New("Must validate one and only one schema (oneOf)"),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
				Context("PatternProperties constraint", func() {
					Context("default props meet all Properties constraints", func() {

						It("returns no errors", func() {

							/* arrange */
							defaultObjectProp1Name := "dummyProp1Name"
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										PatternProperties: map[string]*model.JSONSchema{
											defaultObjectProp1Name: {MinLength: 2},
										},
									},
									Default: map[string]interface{}{
										defaultObjectProp1Name: "dummyProp1Value",
									},
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default props meet some Properties constraints", func() {

						It("returns expected errors", func() {

							/* arrange */
							defaultObjectProp1Name := "dummyProp1Name"
							defaultObjectProp2Name := "dummyProp2Name"

							maxLength := 1
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										PatternProperties: map[string]*model.JSONSchema{
											defaultObjectProp1Name: {MinLength: 2},
											defaultObjectProp2Name: {MaxLength: maxLength},
										},
									},
									Default: map[string]interface{}{
										defaultObjectProp1Name: "dummyProp1Value",
										defaultObjectProp2Name: "dummyProp2Value",
									},
								},
							}

							patterns, err := json.Marshal([]string{defaultObjectProp1Name, defaultObjectProp2Name})
							if err != nil {
								Fail(err.Error())
							}

							expectedErrors := []error{
								fmt.Errorf("String length must be less than or equal to %v", maxLength),
								fmt.Errorf("Property \"%v\" does not match pattern %v", defaultObjectProp2Name, string(patterns)),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							// fix pattern order (maps in go not stable)
							OutOfOrderPatterns, err := json.Marshal([]string{defaultObjectProp2Name, defaultObjectProp1Name})
							if err != nil {
								Fail(err.Error())
							}
							for errIndex, err := range actualErrors {
								actualErrors[errIndex] = errors.New(
									strings.Replace(err.Error(), string(OutOfOrderPatterns), string(patterns), -1),
								)
							}
							Expect(actualErrors).To(ConsistOf(expectedErrors))

						})
					})
					Context("default props meet no Properties constraints", func() {

						It("returns expected errors", func() {

							/* arrange */
							defaultObjectProp1Name := "dummyProp1Name"
							defaultObjectProp2Name := "dummyProp2Name"
							minLength := 100
							maxLength := 1
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										PatternProperties: map[string]*model.JSONSchema{
											defaultObjectProp1Name: {MinLength: minLength},
											defaultObjectProp2Name: {MaxLength: maxLength},
										},
									},
									Default: map[string]interface{}{
										defaultObjectProp1Name: "dummyProp1Value",
										defaultObjectProp2Name: "dummyProp2Value",
									},
								},
							}

							patterns, err := json.Marshal([]string{defaultObjectProp1Name, defaultObjectProp2Name})
							if err != nil {
								Fail(err.Error())
							}

							expectedErrors := []error{
								fmt.Errorf("String length must be greater than or equal to %v", minLength),
								fmt.Errorf("Property \"%v\" does not match pattern %v", defaultObjectProp1Name, string(patterns)),
								fmt.Errorf("String length must be less than or equal to %v", maxLength),
								fmt.Errorf("Property \"%v\" does not match pattern %v", defaultObjectProp2Name, string(patterns)),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							// fix pattern order (maps in go not stable)
							OutOfOrderPatterns, err := json.Marshal([]string{defaultObjectProp2Name, defaultObjectProp1Name})
							if err != nil {
								Fail(err.Error())
							}
							for errIndex, err := range actualErrors {
								actualErrors[errIndex] = errors.New(
									strings.Replace(err.Error(), string(OutOfOrderPatterns), string(patterns), -1),
								)
							}
							Expect(actualErrors).To(ConsistOf(expectedErrors))

						})
					})
				})
				Context("Properties constraint", func() {
					Context("default props meet all Properties constraints", func() {

						It("returns no errors", func() {

							/* arrange */
							defaultObjectProp1Name := "dummyProp1Name"
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										Properties: map[string]*model.JSONSchema{
											defaultObjectProp1Name: {MinLength: 2},
										},
									},
									Default: map[string]interface{}{
										defaultObjectProp1Name: "dummyProp1Value",
									},
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default props meet some Properties constraints", func() {

						It("returns expected errors", func() {

							/* arrange */
							defaultObjectProp1Name := "dummyProp1Name"
							defaultObjectProp2Name := "dummyProp2Name"

							maxLength := 1
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										Properties: map[string]*model.JSONSchema{
											defaultObjectProp1Name: {MinLength: 2},
											defaultObjectProp2Name: {MaxLength: maxLength},
										},
									},
									Default: map[string]interface{}{
										defaultObjectProp1Name: "dummyProp1Value",
										defaultObjectProp2Name: "dummyProp2Value",
									},
								},
							}

							expectedErrors := []error{
								fmt.Errorf("String length must be less than or equal to %v", maxLength),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default props meet no Properties constraints", func() {

						It("returns expected errors", func() {

							/* arrange */
							defaultObjectProp1Name := "dummyProp1Name"
							defaultObjectProp2Name := "dummyProp2Name"

							minLength := 100
							maxLength := 1
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										Properties: map[string]*model.JSONSchema{
											defaultObjectProp1Name: {MinLength: minLength},
											defaultObjectProp2Name: {MaxLength: maxLength},
										},
									},
									Default: map[string]interface{}{
										defaultObjectProp1Name: "dummyProp1Value",
										defaultObjectProp2Name: "dummyProp2Value",
									},
								},
							}

							expectedErrors := []error{
								fmt.Errorf("String length must be greater than or equal to %v", minLength),
								fmt.Errorf("String length must be less than or equal to %v", maxLength),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(ConsistOf(expectedErrors))

						})
					})
				})
				Context("Required constraint", func() {
					Context("default contains all required props", func() {

						It("returns no errors", func() {

							/* arrange */
							defaultObjectProp1Name := "dummyProp1Name"
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										Required: []string{
											defaultObjectProp1Name,
										},
									},
									Default: map[string]interface{}{
										defaultObjectProp1Name: "dummyProp1Value",
									},
								},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default contains no required props", func() {

						It("returns expected errors", func() {

							/* arrange */
							missingPropName := "missingPropName"
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										Required: []string{
											missingPropName,
										},
									},
									Default: map[string]interface{}{
										"dummyProp1Name": "dummyProp1Value",
									},
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"%v is required",
									missingPropName,
								),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
					Context("default contains some required props", func() {

						It("returns expected errors", func() {

							/* arrange */
							defaultObjectProp1Name := "dummyProp1Name"

							missingPropName := "missingPropName"
							providedParam := &model.Param{
								Object: &model.ObjectParam{
									Constraints: &model.ObjectConstraints{
										Required: []string{
											defaultObjectProp1Name,
											missingPropName,
										},
									},
									Default: map[string]interface{}{
										defaultObjectProp1Name: "dummyProp1Value",
										"dummyProp2Name":       "dummyProp2Value",
									},
								},
							}

							expectedErrors := []error{
								fmt.Errorf(
									"%v is required",
									missingPropName,
								),
							}

							/* act */
							actualErrors := objectUnderTest.Validate(nil, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
			})
		})
	})
})
